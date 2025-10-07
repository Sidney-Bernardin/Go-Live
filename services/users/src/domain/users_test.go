package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUsername(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name            string
		username        string
		domainErrorType DomainErrorType
	}{
		{
			name:     "Username",
			username: "foobarbaz",
		},
		{
			name:            "Username Too Short",
			username:        "ab",
			domainErrorType: DomainErrorTypeUsernameInvalid,
		},
		{
			name:            "Username Too Long",
			username:        MustRandomString(34),
			domainErrorType: DomainErrorTypeUsernameInvalid,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Create username.
			username, err := NewUsername(t.Context(), tc.username)

			if tc.domainErrorType == "" {
				assert.Equal(t, tc.username, string(username))
				assert.NoError(t, err)
				return
			}

			var domainErr *DomainError
			if assert.ErrorAs(t, err, &domainErr) {
				assert.Equal(t, tc.domainErrorType, domainErr.Type)
			}

			assert.Equal(t, "", string(username))
		})
	}
}

func TestNewPassword(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name            string
		password        string
		domainErrorType DomainErrorType
	}{
		{
			name:     "Password",
			password: "foobarbaz",
		},
		{
			name:            "Password Too Short",
			password:        "abcdefg",
			domainErrorType: DomainErrorTypePasswordInvalid,
		},
		{
			name:            "Password Too Long",
			password:        MustRandomString(102),
			domainErrorType: DomainErrorTypePasswordInvalid,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Create password.
			password, err := NewPassword(t.Context(), tc.password)

			if tc.domainErrorType == "" {
				assert.Equal(t, tc.password, string(password))
				assert.NoError(t, err)
				return
			}

			var domainErr *DomainError
			if assert.ErrorAs(t, err, &domainErr) {
				assert.Equal(t, tc.domainErrorType, domainErr.Type)
			}

			assert.Equal(t, "", string(password))
		})
	}
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
