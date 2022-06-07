package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  float64(util.RandomBalance()),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func deleteTestAccount(t *testing.T, accountId int64) {
	err := testQueries.DeleteAccount(context.Background(), accountId)
	require.NoError(t, err)
}

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {
	// create account
	account := createRandomAccount(t)

	dbAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)

	require.Equal(t, account.ID, dbAccount.ID)
	require.Equal(t, account.Owner, dbAccount.Owner)
	require.Equal(t, account.Balance, dbAccount.Balance)
	require.Equal(t, account.Currency, dbAccount.Currency)
	require.Equal(t, account.CreatedAt, dbAccount.CreatedAt)

	deleteTestAccount(t, account.ID)
}

func TestUpdateAccount(t *testing.T) {

	account := createRandomAccount(t)

	// using CAD since this cannot be randomly generated
	updateParams := UpdateAccountParams{
		ID:       account.ID,
		Balance:  float64(util.RandomBalance()),
		Currency: "CAD",
	}
	updatedAcc, err := testQueries.UpdateAccount(context.Background(), updateParams)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAcc)

	require.Equal(t, account.ID, updatedAcc.ID)

	require.NotEqual(t, account.Balance, updatedAcc.Balance)
	require.Equal(t, updateParams.Balance, updatedAcc.Balance)

	require.NotEqual(t, account.Currency, updatedAcc.Currency)
	require.Equal(t, updateParams.Currency, updatedAcc.Currency)

	require.Equal(t, account.CreatedAt, updatedAcc.CreatedAt)

	deleteTestAccount(t, account.ID)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	deleteTestAccount(t, account.ID)

	account, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccount(t *testing.T) {

	var accounts []Account

	for i := 0; i < 10; i++ {
		accounts = append(accounts, createRandomAccount(t))
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 3,
	}

	dbAccounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	for _, dbAcc := range dbAccounts {
		require.NotEmpty(t, dbAcc)
	}

	for _, account := range accounts {
		deleteTestAccount(t, account.ID)
	}

}
