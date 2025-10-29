package cache

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
		_, _ = repo.suite(t, 1, 1)
		user := &domain.User{}

		err := repo.InsertUser(t.Context(), user)
		assert.NoError(t, err)
		repo.matchUserSnapshot(t, user.ID)
	})

	t.Run("successful_upsert", func(t *testing.T) {
		users, _ := repo.suite(t, 1, 1)
		user := &domain.User{
			ID: users[0].ID,
		}

		err := repo.InsertUser(t.Context(), user)
		assert.NoError(t, err)
		repo.matchUserSnapshot(t, user.ID, match.Any("id"))
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		users, _ := repo.suite(t, 1, 1)

		user, err := repo.GetUser(t.Context(), users[0].ID)
		assert.NoError(t, err)
		snaps.MatchJSON(t, user, match.Any("ID", "CreatedAt", "UpdatedAt"))
	})

	t.Run("not_found", func(t *testing.T) {
		_, _ = repo.suite(t, 1, 1)

		user, err := repo.GetUser(t.Context(), domain.NewUUID())
		assert.ErrorIs(t, err, service.ErrUserNotFound)
		snaps.MatchJSON(t, user)
	})
}
