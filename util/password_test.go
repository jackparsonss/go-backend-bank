package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPass(t *testing.T) {
	password := RandomString(8)

	hashedPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = Checkpassword(password, hashedPassword)

	require.NoError(t, err)

	wrongPassword := RandomString(8)
	err = Checkpassword(wrongPassword, hashedPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
