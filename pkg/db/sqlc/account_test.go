package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	user := CreateRandomUser(t)

	createArg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), createArg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createArg.Owner, account.Owner)
	require.Equal(t, createArg.Balance, account.Balance)
	require.Equal(t, createArg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreatAccount(t *testing.T) {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}

func TestGetAccount(t *testing.T) {
	// create account
	account1 := createRandomAccount(t)

	arg := GetAccountParams{
		ID:       account1.ID,
		Currency: account1.Currency,
	}
	// get account
	account2, err := testQueries.GetAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.ID, account2.ID)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	//create account
	account1 := createRandomAccount(t)

	//update account
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.ID, account2.ID)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	arg := GetAccountParams{
		ID:       account1.ID,
		Currency: account1.Currency,
	}

	account2, err := testQueries.GetAccount(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestGetAccountByOwner(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := GetAccountByOwnerParams{
		Owner:    account1.Owner,
		Currency: account1.Currency,
	}
	account2, err := testQueries.GetAccountByOwner(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.ID, account2.ID)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}
