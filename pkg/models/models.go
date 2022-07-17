package models

import "time"

type (
	Account struct {
		Id        int64     `db:"id" json:"id"`
		Owner     string    `db:"owner" json:"owner"`
		Currency  string    `db:"currency" json:"currency"`
		Balance   int64     `db:"balance" json:"balance"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		Limit     int64
		Offset    int64
	}

	Entry struct {
		Id        int64     `db:"id" json:"id"`
		AccountID int64     `db:"account_id" json:"account_id"`
		Amount    int64     `db:"amount" json:"amount"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}

	Transfer struct {
		Id            int64     `db:"id" json:"id"`
		FromAccountID int64     `db:"from_account_id" json:"from_account_id"`
		ToAccountID   int64     `db:"to_account_id" json:"to_account_id"`
		Amount        int64     `db:"amount" json:"amount"`
		CreatedAt     time.Time `db:"created_at" json:"created_at"`
	}
)
