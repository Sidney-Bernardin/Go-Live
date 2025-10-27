package cache

import (
	"testing"
	"users/src/domain"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestInsertSession(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		_, _ = repo.suite(t, 1, 1)
		session := &domain.Session{}

		err := repo.InsertSession(t.Context(), session)
		snaps.MatchSnapshot(t, err)
		repo.matchSessionSnapshot(t, session.ID)
	})

	t.Run("successful_upsert", func(t *testing.T) {
		_, sessions := repo.suite(t, 1, 1)
		session := &domain.Session{
			ID: sessions[0].ID,
		}

		err := repo.InsertSession(t.Context(), session)
		snaps.MatchSnapshot(t, err)
		repo.matchSessionSnapshot(t, session.ID, match.Any("id"))
	})
}

func TestGetSession(t *testing.T) {
	t.Parallel()
	repo := suite(t)

	t.Run("successful", func(t *testing.T) {
		_, sessions := repo.suite(t, 1, 1)

		session, err := repo.GetSession(t.Context(), sessions[0].ID)
		snaps.MatchSnapshot(t, err)
		snaps.MatchJSON(t, session, match.Any("ID", "UserID", "CSRFToken"))
	})

	t.Run("not_found", func(t *testing.T) {
		_, _ = repo.suite(t, 1, 1)

		session, err := repo.GetSession(t.Context(), domain.NewUUID())
		snaps.MatchSnapshot(t, err)
		snaps.MatchJSON(t, session)
	})
}
