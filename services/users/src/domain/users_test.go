package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUsername(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		u := "foobarbaz"

		username, err := NewUsername(u)
		assert.Equal(t, u, string(username))
		assert.NoError(t, err)
	})

	t.Run("too_short", func(t *testing.T) {
		u := ""

		username, err := NewUsername(u)
		assert.Equal(t, "", string(username))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypeUsernameInvalid, domainErr.Type)
		}
	})

	t.Run("too_long", func(t *testing.T) {
		u := string(make([]byte, 34))

		username, err := NewUsername(u)
		assert.Equal(t, "", string(username))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypeUsernameInvalid, domainErr.Type)
		}
	})
}

func TestNewPassword(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		pw := "foobarbaz"

		password, err := NewPassword(pw)
		assert.Equal(t, pw, string(password))
		assert.NoError(t, err)
	})

	t.Run("too_short", func(t *testing.T) {
		pw := ""

		password, err := NewPassword(pw)
		assert.Equal(t, "", string(password))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypePasswordInvalid, domainErr.Type)
		}
	})

	t.Run("too_long", func(t *testing.T) {
		pw := string(make([]byte, 102))

		password, err := NewPassword(pw)
		assert.Equal(t, "", string(password))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypePasswordInvalid, domainErr.Type)
		}
	})
}

func TestNewPasswordSalt(t *testing.T) {
	t.Parallel()

	t.Run("length", func(t *testing.T) {
		assert.Equal(t, len(NewPasswordSalt()), 32)
	})
}

func TestNewPasswordHash(t *testing.T) {
	t.Parallel()

	t.Run("bcrypt_comparison", func(t *testing.T) {
		password, err := NewPassword("foobarbaz")
		require.NoError(t, err)

		passwordSalt := NewPasswordSalt()

		passwordHash, err := NewPasswordHash(password, passwordSalt)
		require.NoError(t, err)
		require.NoError(t, bcrypt.CompareHashAndPassword(passwordHash, []byte(string(password)+string(passwordSalt))))
	})
}
