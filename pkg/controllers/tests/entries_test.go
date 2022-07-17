package controllers

import (
	"context"
	accountController "simplebank/pkg/controllers/account"
	entryController "simplebank/pkg/controllers/entry"
	"simplebank/pkg/models"
	"simplebank/pkg/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account *models.Account) *models.Entry {

	entry1 := entryController.CreateEntryParams{
		AccountID: account.Id,
		Amount:    util.RandomMoney(),
	}

	entry2, err := entryController.CreateEntry(context.Background(), DB, entry1)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.NotZero(t, entry2.Id)
	require.NotZero(t, entry2.CreatedAt)

	return entry2
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)

	entry1 := createRandomEntry(t, account)
	require.NotEmpty(t, entry1)

	require.Equal(t, account.Id, entry1.AccountID)

	args := accountController.UpdateAccountParams{
		Id:      account.Id,
		Balance: account.Balance + entry1.Amount,
	}
	account1, err := accountController.UpdateAccount(context.Background(), DB, args)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, entry1.AccountID, account1.Id)
	require.Equal(t, account1.Balance, (account.Balance + entry1.Amount))
}

func TestGetEntryByID(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	entry1, err := entryController.GetEntryByID(context.Background(), DB, entry.Id)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	require.Equal(t, entry.Id, entry1.Id)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.Amount, entry1.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry1.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry := createRandomEntry(t, account)

	entry1, err := entryController.DeleteEntry(context.Background(), DB, entry.Id)
	require.NoError(t, err)
	require.NotEmpty(t, entry1)

	require.Equal(t, entry.Id, entry1)

	args := accountController.UpdateAccountParams{
		Id:      account.Id,
		Balance: account.Balance - entry.Amount,
	}
	account1, err := accountController.UpdateAccount(context.Background(), DB, args)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, account1.Balance, (account.Balance - entry.Amount))
}
