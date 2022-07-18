package controller

import (
	"context"
	"database/sql"
	accountController "simplebank/pkg/controllers/account"
	entryController "simplebank/pkg/controllers/entry"
	"simplebank/pkg/models"
	"sync"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type (
	TransferTxParams struct {
		FromAccountID int64 `json:"from_account_id"`
		ToAccountID   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}

	TransferTxResult struct {
		Transfer    *models.Transfer
		FromAccount *models.Account
		ToAccount   *models.Account
		EntryFrom   *models.Entry
		EntryTo     *models.Entry
	}
)

var mutex sync.Mutex

func execTx(ctx context.Context, db *sqlx.DB, fn func() error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed begin transaction")
	}

	err = fn()
	if err != nil {
		if rtx := tx.Rollback(); rtx != nil {
			return errors.Wrapf(err, "failed tx err: %v, rtx err: %v", err, rtx)
		}
		return errors.Wrap(err, "failed transaction")
	}

	return tx.Commit()
}

func TransferTx(ctx context.Context, db *sqlx.DB, args TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult

	err := execTx(ctx, db, func() error {
		var err error

		result.Transfer, err = CreateTransfer(ctx, db, args)
		if err != nil {
			return err
		}

		result.EntryFrom, err = entryController.CreateEntry(ctx, db, entryController.CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.EntryTo, err = entryController.CreateEntry(ctx, db, entryController.CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		// Get account & Update account's balance
		mutex.Lock()
		account1, err := accountController.GetAccountByIDForUpdate(context.Background(), db, args.FromAccountID)
		if err != nil {
			return err
		}
		result.FromAccount, err = accountController.UpdateAccount(context.Background(), db, accountController.UpdateAccountParams{
			Id:      args.FromAccountID,
			Balance: account1.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		account2, err := accountController.GetAccountByIDForUpdate(context.Background(), db, args.ToAccountID)
		if err != nil {
			return err
		}
		result.ToAccount, err = accountController.UpdateAccount(context.Background(), db, accountController.UpdateAccountParams{
			Id:      args.ToAccountID,
			Balance: account2.Balance + args.Amount,
		})
		if err != nil {
			return err
		}
		mutex.Unlock()

		return nil
	})
	if err != nil {
		return &result, errors.Wrap(err, "failed execTx")
	}

	return &result, nil
}

func CreateTransfer(ctx context.Context, db *sqlx.DB, args TransferTxParams) (*models.Transfer, error) {
	query := `INSERT INTO transfers ("from_account_id", "to_account_id", "amount") VALUES ($1, $2, $3) RETURNING *`

	var transfer models.Transfer
	row, err := db.QueryContext(ctx, query, args.FromAccountID, args.ToAccountID, args.Amount)
	if err == sql.ErrNoRows {
		return &transfer, nil
	}
	if err != nil {
		return &transfer, errors.Wrap(err, "failed insert")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&transfer.Id, &transfer.FromAccountID, &transfer.ToAccountID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return &transfer, errors.Wrap(err, "failed scan")
		}
	}

	// log.Println("Transfer created")
	return &transfer, nil
}

func GetTransferByID(ctx context.Context, db *sqlx.DB, id int64) (*models.Transfer, error) {
	query := `SELECT * FROM transfers WHERE id = $1 LIMIT 1`

	var res models.Transfer
	row, err := db.QueryContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed retrieving the row")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.FromAccountID, &res.ToAccountID, &res.Amount, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}
	}

	return &res, nil
}
