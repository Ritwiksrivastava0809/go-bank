package db

import (
	"context"
	"testing"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	password := utils.RandomPassword()
	hashedPassword, err := utils.HashPasswordArgon2(password)
	require.NoError(t, err)

	user1 := CreateUserParams{
		Username:          utils.RandomUsername(),
		Email:             utils.RandomEmail(),
		FullName:          utils.RandomOwner(),
		HashedPassword:    hashedPassword,
		PasswordChangedAt: time.Now(),
	}

	user, err := testQueries.CreateUser(context.Background(), user1)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestCreateUSer(t *testing.T) {

	password := utils.RandomPassword()
	hashedPassword, err := utils.HashPasswordArgon2(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:          utils.RandomUsername(),
		Email:             utils.RandomEmail(),
		HashedPassword:    hashedPassword,
		FullName:          utils.RandomOwner(),
		PasswordChangedAt: time.Now(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.WithinDuration(t, arg.PasswordChangedAt, user.CreatedAt, time.Second)

	require.NotZero(t, user.CreatedAt)

}

func TestGetUserByEmail(t *testing.T) {

	user1 := CreateRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.NotZero(t, user2.CreatedAt)
}

func TestGetUserByUsername(t *testing.T) {

	user1 := CreateRandomUser(t)

	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.NotZero(t, user2.CreatedAt)
}
