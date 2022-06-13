package db

import (
	"context"
	"testing"

	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashPass, err := util.HashPassword(password)

	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hashPass,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	// create account
	user := createRandomUser(t)

	dbUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, dbUser)

	require.Equal(t, user.Username, dbUser.Username)
	require.NotEmpty(t, dbUser.HashedPassword)
	require.NotEmpty(t, dbUser.FullName)
	require.Equal(t, user.Email, dbUser.Email)
	require.Equal(t, user.CreatedAt, dbUser.CreatedAt)
	require.NotEmpty(t, dbUser.PasswordChangedAt)

}
