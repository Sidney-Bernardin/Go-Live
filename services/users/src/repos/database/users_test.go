package database

import (
	"testing"
	"users/src/domain"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {

	ctx := t.Context()
	repo := newTest(t)

	tt := []struct {
		name    string
		oldUser *domain.User
		newUser *domain.User
		err     error
	}{
		{
			name:    "Insert",
			newUser: &domain.User{ID: domain.NewUUID(), PasswordHash: []byte(``)},
		},
		{
			name:    "Username Taken",
			oldUser: &domain.User{ID: domain.NewUUID(), Username: "foo", Email: "foo", PasswordHash: []byte(``)},
			newUser: &domain.User{ID: domain.NewUUID(), Username: "foo", Email: "bar", PasswordHash: []byte(``)},
			err:     domain.ErrUsernameTaken,
		},
		{
			name:    "Email Taken",
			oldUser: &domain.User{ID: domain.NewUUID(), Username: "foo", Email: "foo", PasswordHash: []byte(``)},
			newUser: &domain.User{ID: domain.NewUUID(), Username: "bar", Email: "foo", PasswordHash: []byte(``)},
			err:     domain.ErrEmailTaken,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			t.Cleanup(func() {
				// Clear the users table.
				_, err := repo.pool.Exec(ctx, `DELETE FROM users`)
				require.NoError(t, err)
			})

			if tc.oldUser != nil {

				// Insert old user.
				_, err := repo.pool.Exec(ctx,
					`
						INSERT INTO users (id, username, email, password_hash, password_salt)
						VALUES ($1, $2, $3, $4, $5)
					`,
					tc.oldUser.ID,
					tc.oldUser.Username,
					tc.oldUser.Email,
					tc.oldUser.PasswordHash,
					tc.oldUser.PasswordSalt)
				require.NoError(t, err)
			}

			// Do test with new user.
			err := repo.InsertUser(ctx, tc.newUser)
			if !assert.ErrorIs(t, err, tc.err) || err != nil {
				return
			}

			// Get the new user.
			q := `SELECT id, username, email, password_hash, password_salt FROM users WHERE id = $1`
			rows, err := repo.pool.Query(ctx, q, tc.newUser.ID)
			require.NoError(t, err)

			// Decode the new user.
			u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dbUser])
			require.NoError(t, err)

			assert.Equal(t, tc.newUser.ID, u.ID)
			assert.Equal(t, tc.newUser.Username, u.Username)
			assert.Equal(t, tc.newUser.Email, u.Email)
			assert.Equal(t, tc.newUser.PasswordHash, u.PasswordHash)
			assert.Equal(t, tc.newUser.PasswordSalt, u.PasswordSalt)
		})
	}
}
