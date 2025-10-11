package database

import (
	"testing"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertUser(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	repo := newTestDatabase(t)

	// Insert a new dummy user.
	dummyUser := newTestUser(t)
	_, err := repo.pool.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash, password_salt) VALUES ($1, $2, $3, $4, $5)`,
		dummyUser.ID, dummyUser.Username, dummyUser.Email, dummyUser.PasswordHash, dummyUser.PasswordSalt)
	require.NoError(t, err)

	t.Run("work", func(t *testing.T) {
		user := newTestUser(t)

		err := repo.InsertUser(t.Context(), user)
		assert.NoError(t, err)
		repo.assertUser(t, user.ID, user)
	})

	t.Run("username_conflict", func(t *testing.T) {
		user := newTestUser(t)
		user.Username = dummyUser.Username

		err := repo.InsertUser(t.Context(), user)
		assert.Equal(t, service.ErrUsernameTaken, err)
		repo.assertUser(t, user.ID, nil)
	})

	t.Run("email_conflict", func(t *testing.T) {
		user := newTestUser(t)
		user.Email = dummyUser.Email

		err := repo.InsertUser(t.Context(), user)
		assert.Equal(t, service.ErrEmailTaken, err)
		repo.assertUser(t, user.ID, nil)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	db := newTestDatabase(t)

	// Insert a new dummy user.
	dummyUser := newTestUser(t)
	_, err := db.pool.Exec(ctx,
		`INSERT INTO users (id, username, email, password_hash, password_salt) VALUES ($1, $2, $3, $4, $5)`,
		dummyUser.ID, dummyUser.Username, dummyUser.Email, dummyUser.PasswordHash, dummyUser.PasswordSalt)
	require.NoError(t, err)

	t.Run("work", func(t *testing.T) {
		user, err := db.GetUserByID(t.Context(), dummyUser.ID)
		assert.NoError(t, err)
		db.assertUser(t, dummyUser.ID, user)
	})

	t.Run("not_found", func(t *testing.T) {
		user := newTestUser(t)
		user.Username = dummyUser.Username

		user, err := db.GetUserByID(t.Context(), domain.NewUUID())
		assert.Equal(t, service.ErrUserNotFound, err)
	})
}
