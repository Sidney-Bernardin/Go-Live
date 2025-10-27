package database

import (
	"testing"
	"users/src/domain"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestInsertUser(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		_ = repo.suite(t, 0)
		err := repo.InsertUser(t.Context(), &domain.User{
			PasswordHash: []byte(``),
		})
		snaps.MatchSnapshot(t, err)
	})

	t.Run("username_taken", func(t *testing.T) {
		users := repo.suite(t, 1)
		err := repo.InsertUser(t.Context(), &domain.User{
			Username:     users[0].Username,
			PasswordHash: []byte(``),
		})
		snaps.MatchSnapshot(t, err)
	})

	t.Run("email_taken", func(t *testing.T) {
		users := repo.suite(t, 1)
		err := repo.InsertUser(t.Context(), &domain.User{
			Email:        users[0].Email,
			PasswordHash: []byte(``),
		})
		snaps.MatchSnapshot(t, err)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		users := repo.suite(t, 1)
		user, err := repo.GetUserByID(t.Context(), users[0].ID)

		snaps.MatchSnapshot(t, err)
		snaps.MatchJSON(t, user, match.Any("ID", "CreatedAt", "UpdatedAt"))
	})

	t.Run("not_found", func(t *testing.T) {
		_ = repo.suite(t, 1)
		user, err := repo.GetUserByID(t.Context(), domain.NewUUID())

		snaps.MatchSnapshot(t, err)
		snaps.MatchJSON(t, user)
	})
}
