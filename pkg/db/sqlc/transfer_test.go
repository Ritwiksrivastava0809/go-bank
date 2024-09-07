package db

import (
	"context"
	"testing"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransaction(t *testing.T, account1 int64, account2 int64) Transfer {
	arg := InsertTransactionParams{
		FromAccountID: account1,
		ToAccountID:   account2,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.InsertTransaction(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestInsertTransaction(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := InsertTransactionParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.InsertTransaction(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransactionByID(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer1 := CreateRandomTransaction(t, account1.ID, account2.ID)

	transfer2, err := testQueries.GetTransactionByID(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

	require.NotZero(t, transfer2.FromAccountID)
	require.NotZero(t, transfer2.ToAccountID)
	require.NotZero(t, transfer2.CreatedAt)
}

func TestGetTransactionHistoryByAccountID(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomTransaction(t, account1.ID, account2.ID)
	}

	arg := GetTransactionHistoryByAccountIDParams{
		FromAccountID: account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.GetTransactionHistoryByAccountID(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.FromAccountID)
		require.NotZero(t, transfer.ToAccountID)
		require.NotZero(t, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)

	}
}

func TestDeleteTransactionByID(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer1 := CreateRandomTransaction(t, account1.ID, account2.ID)

	err := testQueries.DeleteTransactionByID(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransactionByID(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)

}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomTransaction(t, account1.ID, account2.ID)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.FromAccountID)
		require.NotZero(t, transfer.ToAccountID)
		require.NotZero(t, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
	}

}
