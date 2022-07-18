package controllers

import (
	"context"
	"database/sql"
	"log"
	"simplebank/pkg/models"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type (
	CreateEntryParams struct {
		AccountID int64 `json:"account_id"`
		Amount    int64 `json:"amount"`
	}

	ListEntryParams struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}

	UpdateEntryParams struct {
		Id     int64 `json:"id"`
		Amount int64 `json:"amount"`
	}
)

func CreateEntry(ctx context.Context, db *sqlx.DB, entry CreateEntryParams) (*models.Entry, error) {
	query := `INSERT INTO entries ("account_id", "amount") VALUES ($1, $2) RETURNING *`

	var res models.Entry
	row, err := db.QueryContext(ctx, query, entry.AccountID, entry.Amount)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed insert")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.AccountID, &res.Amount, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}
	}

	// log.Println("Entry created")
	return &res, nil
}

func GetEntryByID(ctx context.Context, db *sqlx.DB, id int64) (*models.Entry, error) {
	query := `SELECT * FROM entries WHERE id = $1 LIMIT 1`

	var res models.Entry
	row, err := db.QueryContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed retrieving the row")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.AccountID, &res.Amount, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}
	}

	return &res, nil
}

func GetEntryAll(ctx context.Context, db *sqlx.DB, args ListEntryParams) (*[]models.Entry, error) {
	query := `SELECT * FROM entries ORDER BY id LIMIT $1 OFFSET $2`

	var res []models.Entry
	rows, err := db.QueryContext(ctx, query, args.Limit, args.Offset)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed retrieving the row")
	}

	defer rows.Close()
	for rows.Next() {
		var entry models.Entry
		err := rows.Scan(&entry.Id, &entry.AccountID, &entry.Amount, &entry.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}

		res = append(res, entry)
	}

	return &res, nil
}

func UpdateEntry(ctx context.Context, db *sqlx.DB, args UpdateEntryParams) (*models.Entry, error) {
	query := `UPDATE entries SET amount = $1 WHERE id = $2 RETURNING *`

	var res models.Entry
	row, err := db.QueryContext(ctx, query, args.Amount, args.Id)
	if err == sql.ErrNoRows {
		return &res, nil
	}
	if err != nil {
		return &res, errors.Wrap(err, "failed update")
	}

	defer row.Close()
	for row.Next() {
		err := row.Scan(&res.Id, &res.AccountID, &res.Amount, &res.CreatedAt)
		if err != nil {
			return &res, errors.Wrap(err, "failed scan")
		}
	}

	log.Println("Entry Updated")
	return &res, nil
}

func DeleteEntry(ctx context.Context, db *sqlx.DB, id int64) (int64, error) {
	query := `DELETE FROM entries WHERE id = $1 RETURNING id`

	var res int64
	row, err := db.QueryContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return res, errors.Wrap(err, "any record retrieved")
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

	log.Println("Entry Deleted")
	return res, nil
}
