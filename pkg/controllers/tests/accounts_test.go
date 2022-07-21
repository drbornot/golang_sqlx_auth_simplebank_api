package controllers

import (
	"context"
	accountController "simplebank/pkg/controllers/account"
	"simplebank/pkg/models"
	"simplebank/pkg/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *models.Account {
	account := accountController.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}
	res, err := accountController.CreateAccount(context.Background(), DB, account)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, res.Owner, account.Owner)
	require.Equal(t, res.Balance, account.Balance)
	require.Equal(t, res.Currency, account.Currency)

	require.NotZero(t, res.Id)
	require.NotZero(t, res.CreatedAt)

	return res
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountByID(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := accountController.GetAccountByID(context.Background(), DB, account1.Id)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Id, account2.Id)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestGetAccountAll(t *testing.T) {
	args := accountController.ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	accounts, err := accountController.GetAccountAll(context.Background(), DB, args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	args := accountController.UpdateAccountParams{
		Id:      account1.Id,
		Balance: util.RandomMoney(),
	}

	account2, err := accountController.UpdateAccount(context.Background(), DB, args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Id, account2.Id)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	id, err := accountController.DeleteAccout(context.Background(), DB, account.Id)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	require.Equal(t, account.Id, id)
}
