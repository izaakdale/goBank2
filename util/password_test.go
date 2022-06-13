package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashing(t *testing.T) {
	password := RandomString(6)

	// test hashing
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	// test verify
	err = VerifyPassword(password, hash1)
	require.NoError(t, err)

	// test incorrect password
	incorrectPassword := "nope"
	err = VerifyPassword(incorrectPassword, hash1)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// test same pass create two hashes
	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)
}
