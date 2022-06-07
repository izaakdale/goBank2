package db

import (
	"context"
	"testing"

	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, account Account) (Entry, CreateEntryParams) {

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    float64(util.RandomInt(1, 200)),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)

	return entry, args
}

func TestCreateEntry(t *testing.T) {

	account := createRandomAccount(t)
	entry, args := createTestEntry(t, account)

	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	deleteTestAccount(t, account.ID)
}

func TestGetEntry(t *testing.T) {

	account := createRandomAccount(t)
	entry, _ := createTestEntry(t, account)

	dbEntry, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.Equal(t, entry.AccountID, dbEntry.AccountID)
	require.Equal(t, entry.Amount, dbEntry.Amount)
	require.Equal(t, entry.ID, dbEntry.ID)
	require.Equal(t, entry.CreatedAt, dbEntry.CreatedAt)

	deleteTestAccount(t, account.ID)
}

func TestListEntries(t *testing.T) {

	account := createRandomAccount(t)
	var entries []Entry

	for i := 0; i < 10; i++ {
		entry, _ := createTestEntry(t, account)
		entries = append(entries, entry)
	}

	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    3,
	}
	dbEntries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)

	if len(dbEntries) < 1 {
		t.Error("Failed to get any rows from db")
	}

	for _, dbEntry := range dbEntries {
		require.NotEmpty(t, dbEntry)
	}

	deleteTestAccount(t, account.ID)
}
