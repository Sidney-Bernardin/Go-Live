package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUsername(t *testing.T) {
	t.Parallel()

	t.Run("work", func(t *testing.T) {
		inUsername := "foobarbaz"
		outUsername, err := NewUsername(t.Context(), inUsername)

		assert.Equal(t, inUsername, string(outUsername))
		assert.NoError(t, err)
	})

	t.Run("too_short", func(t *testing.T) {
		inUsername := ""
		outUsername, err := NewUsername(t.Context(), inUsername)

		assert.Equal(t, "", string(outUsername))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypeUsernameInvalid, domainErr.Type)
		}
	})

	t.Run("too_long", func(t *testing.T) {
		inUsername := MustRandomString(34)
		outUsername, err := NewUsername(t.Context(), inUsername)

		assert.Equal(t, "", string(outUsername))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypeUsernameInvalid, domainErr.Type)
		}
	})
}

func TestNewPassword(t *testing.T) {
	t.Parallel()

	t.Run("work", func(t *testing.T) {
		inPassword := "foobarbaz"
		outPassword, err := NewPassword(t.Context(), inPassword)

		assert.Equal(t, inPassword, string(outPassword))
		assert.NoError(t, err)
	})

	t.Run("too_short", func(t *testing.T) {
		inPassword := ""
		outPassword, err := NewPassword(t.Context(), inPassword)

		assert.Equal(t, "", string(outPassword))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypePasswordInvalid, domainErr.Type)
		}
	})

	t.Run("too_long", func(t *testing.T) {
		inPassword := MustRandomString(102)
		outPassword, err := NewPassword(t.Context(), inPassword)

		assert.Equal(t, "", string(outPassword))
		if domainErr := (&DomainError{}); assert.ErrorAs(t, err, &domainErr) {
			assert.Equal(t, DomainErrorTypePasswordInvalid, domainErr.Type)
		}
	})
}

func TestNewPasswordSalt(t *testing.T) {
	t.Parallel()

	t.Run("Length", func(t *testing.T) {
		assert.Equal(t, len(NewPasswordSalt()), 32)
	})
}

func TestNewPasswordHash(t *testing.T) {
	t.Parallel()

	t.Run("Bcrypt Comparison", func(t *testing.T) {
		ctx := t.Context()

		// Create password.
		password, err := NewPassword(ctx, "foobarbaz")
		require.NoError(t, err)

		// Create password-salt.
		passwordSalt := NewPasswordSalt()

		// Create password-hash.
		passwordHash, err := NewPasswordHash(ctx, password, passwordSalt)
		require.NoError(t, err)

		// Assert bcrypt comparison.
		err = bcrypt.CompareHashAndPassword(passwordHash, []byte(string(password)+string(passwordSalt)))
		assert.NoError(t, err)
	})
}
