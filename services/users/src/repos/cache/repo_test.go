package cache

import (
	"fmt"
	"testing"
	"time"
	"users/src"
	"users/src/domain"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func suite(t *testing.T) *repository {
	t.Helper()
	ctx := t.Context()

	redisContainer, err := redis.Run(ctx, "redis:8.2.2")
	testcontainers.CleanupContainer(t, redisContainer)
	require.NoError(t, err)

	redisAddr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := New(ctx, &src.Config{
		SessionDuration: time.Hour,
		RedisAddr:       redisAddr[8:],
	})
	require.NoError(t, err)

	return repo.(*repository)
}

func (repo *repository) suite(t *testing.T, users, sessions int) ([]*domain.User, []*domain.Session) {
	t.Helper()
	ctx := t.Context()

	err := repo.client.FlushAll(ctx).Err()
	require.NoError(t, err)

	var dummyUsers []*domain.User
	for i := range users {
		u := repoUser{
			ID:       uuid.New(),
			Username: fmt.Sprintf("username%v", i),
			Email:    fmt.Sprintf("email%v", i),
		}

		err = repo.client.JSONSet(ctx, fmt.Sprintf("user:%s", u.ID), ".", u).Err()
		require.NoError(t, err)

		dummyUsers = append(dummyUsers, u.domainify())
	}

	var dummySessions []*domain.Session
	for range sessions {
		s := repoSession{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			CSRFToken: string(domain.NewCSRFToken()),
		}

		err = repo.client.JSONSet(ctx, fmt.Sprintf("session:%s", s.ID), ".", s).Err()
		require.NoError(t, err)

		dummySessions = append(dummySessions, s.domainify())
	}

	return dummyUsers, dummySessions
}

func (repo *repository) matchUserSnapshot(t *testing.T, userID domain.UUID, matchers ...match.JSONMatcher) {
	t.Helper()
	user, err := repo.client.JSONGet(t.Context(), fmt.Sprintf("user:%s", userID), ".").Result()
	snaps.MatchSnapshot(t, err)
	snaps.MatchJSON(t, user, matchers...)
}

func (repo *repository) matchSessionSnapshot(t *testing.T, sessionID domain.UUID, matchers ...match.JSONMatcher) {
	t.Helper()
	session, err := repo.client.JSONGet(t.Context(), fmt.Sprintf("session:%s", sessionID), ".").Result()
	snaps.MatchSnapshot(t, err)
	snaps.MatchJSON(t, session, matchers...)
}
