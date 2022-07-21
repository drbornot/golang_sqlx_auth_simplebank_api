package controllers

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"simplebank/pkg/models"
)

type (
	ListAccountParams struct {
		Limit  int32 `db:"limit" json:"limit"`
		Offset int32 `db:"offset" json:"offset"`
	}

	CreateAccountParams struct {
		Owner    string `db:"owner" json:"owner"`
		Currency string `db:"currency" json:"currency"`
		Balance  int64  `db:"balance" json:"balance"`
	}

	UpdateAccountParams struct {
		Id      int64 `db:"id" json:"id"`
		Balance int64 `db:"balance" json:"balance"`
	}
)

func CreateAccount(ctx context.Context, db *sqlx.DB, account CreateAccountParams) (*models.Account, error) {
	query := `INSERT INTO accounts ("owner", "currency", "balance") VALUES ($1, $2, $3) RETURNING *`

	var res models.Account
	row, err := db.QueryContext(ctx, query, account.Owner, account.Currency, account.Balance)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed insert")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed row scan")
		}
	}

	// log.Println("Account created")
	return &res, nil
}

func GetAccountByID(ctx context.Context, db *sqlx.DB, id int64) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1`

	var res models.Account
	err := db.QueryRowContext(ctx, query, id).Scan(&res.Id, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
	if err == sql.ErrNoRows {
		return &res, errors.Wrap(err, "row not found")
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed retrieving the row")
	}

	return &res, nil
}

func GetAccountByIDForUpdate(ctx context.Context, db *sqlx.DB, id int64) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1 FOR UPDATE;`

	var res models.Account
	row, err := db.QueryContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed retrieving the row")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed row scan")
		}
	}

	return &res, nil
}

func GetAccountAll(ctx context.Context, db *sqlx.DB, arg ListAccountParams) ([]models.Account, error) {
	query := `SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2`

	var res []models.Account
	rows, err := db.QueryContext(ctx, query, arg.Limit, arg.Offset)
	if err == sql.ErrNoRows {
		return res, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed insert")
	}

	defer rows.Close()
	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.Id, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed rows scan")
		}

		res = append(res, account)
	}

	return res, nil
}

func UpdateAccount(ctx context.Context, db *sqlx.DB, arg UpdateAccountParams) (*models.Account, error) {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING *`

	var res models.Account
	row, err := db.QueryContext(ctx, query, arg.Balance, arg.Id)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed update")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}
	}

	// log.Println("Account Updated")
	return &res, nil
}

func DeleteAccout(ctx context.Context, db *sqlx.DB, id int64) (int64, error) {
	query := `DELETE FROM accounts WHERE id = $1 RETURNING id`

	var res int64
	row, err := db.QueryContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return res, nil
	}
	if err != nil {
		return res, errors.Wrap(err, "failed delete")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res)
		if err != nil {
			return res, errors.Wrap(err, "failed scan")
		}
	}

	log.Println("Account Deleted")
	return res, nil
}
