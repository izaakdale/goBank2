package db

import (
	"context"
	"testing"

	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T, account1, account2 Account) (Transfer, CreateTransferParams) {

	args := CreateTransferParams{
		ToAccountID:   account1.ID,
		FromAccountID: account2.ID,
		Amount:        float64(util.RandomBalance()),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)

	return transfer, args
}

func TestCreateTransfer(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer, args := createTestTransfer(t, account1, account2)

	require.Equal(t, args.Amount, transfer.Amount)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.NotEmpty(t, transfer.ID)
	require.NotEmpty(t, transfer.CreatedAt)

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}

func TestGetTransfer(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer, _ := createTestTransfer(t, account1, account2)

	dbTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	require.Equal(t, transfer.ID, dbTransfer.ID)
	require.Equal(t, transfer.FromAccountID, dbTransfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, dbTransfer.ToAccountID)
	require.Equal(t, transfer.Amount, dbTransfer.Amount)
	require.NotEmpty(t, dbTransfer.ID)
	require.NotEmpty(t, dbTransfer.CreatedAt)

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}

func TestListTransfers(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	var transfers []Transfer

	for i := 0; i < 10; i++ {
		transfer, _ := createTestTransfer(t, account1, account2)
		transfers = append(transfers, transfer)
	}

	args := ListTransfersParams{
		ToAccountID:   account1.ID,
		FromAccountID: account2.ID,
		Limit:         5,
		Offset:        3,
	}

	dbTransfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)

	if len(dbTransfers) < 1 {
		t.Error("Failed to get rows from db")
	}

	for _, dbTransfer := range dbTransfers {
		require.NotEmpty(t, dbTransfer)
	}

	deleteTestAccount(t, account1.ID)
	deleteTestAccount(t, account2.ID)
}
