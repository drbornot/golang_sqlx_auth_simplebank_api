package controllers

import (
	"context"
	accountController "simplebank/pkg/controllers/account"
	entryController "simplebank/pkg/controllers/entry"
	transferController "simplebank/pkg/controllers/transfer"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan transferController.TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := transferController.TransferTx(context.Background(), DB, transferController.TransferTxParams{
				FromAccountID: account1.Id,
				ToAccountID:   account2.Id,
				Amount:        amount,
			})
			errs <- err
			results <- *result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.Id, transfer.FromAccountID)
		require.Equal(t, account2.Id, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.Id)
		require.NotZero(t, transfer.CreatedAt)

		transfer1, err := transferController.GetTransferByID(context.Background(), DB, transfer.Id)
		require.NoError(t, err)
		require.NotEmpty(t, transfer1)

		require.Equal(t, transfer.Id, transfer1.Id)

		// check entries
		entryForm := result.EntryFrom
		require.NotEmpty(t, entryForm)
		require.Equal(t, account1.Id, entryForm.AccountID)
		require.Equal(t, -amount, entryForm.Amount)
		require.NotZero(t, entryForm.Id)
		require.NotZero(t, entryForm.CreatedAt)

		entryForm1, err := entryController.GetEntryByID(context.Background(), DB, entryForm.Id)
		require.NoError(t, err)
		require.NotEmpty(t, entryForm1)
		require.Equal(t, entryForm.Id, entryForm1.Id)

		entryTo := result.EntryTo
		require.NotEmpty(t, entryTo)
		require.Equal(t, account2.Id, entryTo.AccountID)
		require.Equal(t, amount, entryTo.Amount)
		require.NotZero(t, entryTo.Id)
		require.NotZero(t, entryTo.CreatedAt)

		entryTo1, err := entryController.GetEntryByID(context.Background(), DB, entryTo.Id)
		require.NoError(t, err)
		require.NotEmpty(t, entryTo1)
		require.Equal(t, entryTo.Id, entryTo1.Id)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.Id, fromAccount.Id)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.Id, toAccount.Id)

		// check account's balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
	}

	updatedAccount1, err := accountController.GetAccountByID(context.Background(), DB, account1.Id)
	require.NoError(t, err)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)

	updatedAccount2, err := accountController.GetAccountByID(context.Background(), DB, account2.Id)
	require.NoError(t, err)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
