package database

import (
	"testing"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		_ = repo.suite(t, 0)
		err := repo.InsertUser(t.Context(), &domain.User{
			PasswordHash: []byte(``),
		})
		assert.NoError(t, err)
	})

	t.Run("username_taken", func(t *testing.T) {
		users := repo.suite(t, 1)
		err := repo.InsertUser(t.Context(), &domain.User{
			Username:     users[0].Username,
			PasswordHash: []byte(``),
		})
		assert.ErrorIs(t, err, service.ErrUsernameTaken)
	})

	t.Run("email_taken", func(t *testing.T) {
		users := repo.suite(t, 1)
		err := repo.InsertUser(t.Context(), &domain.User{
			Email:        users[0].Email,
			PasswordHash: []byte(``),
		})
		assert.ErrorIs(t, err, service.ErrEmailTaken)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		users := repo.suite(t, 1)
		user, err := repo.GetUserByID(t.Context(), users[0].ID)

		assert.NoError(t, err)
		snaps.MatchJSON(t, user, match.Any("ID", "CreatedAt", "UpdatedAt"))
	})

	t.Run("not_found", func(t *testing.T) {
		_ = repo.suite(t, 1)
		user, err := repo.GetUserByID(t.Context(), domain.NewUUID())

		assert.ErrorIs(t, err, service.ErrUserNotFound)
		snaps.MatchJSON(t, user)
	})
}

func TestGetUserByUsername(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		users := repo.suite(t, 1)
		user, err := repo.GetUserByUsername(t.Context(), users[0].Username)

		assert.NoError(t, err)
		snaps.MatchJSON(t, user, match.Any("ID", "CreatedAt", "UpdatedAt"))
	})

	t.Run("not_found", func(t *testing.T) {
		_ = repo.suite(t, 1)
		user, err := repo.GetUserByUsername(t.Context(), domain.Username(""))

		assert.ErrorIs(t, err, service.ErrUserNotFound)
		snaps.MatchJSON(t, user)
	})
}
