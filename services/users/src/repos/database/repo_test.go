package database

import (
	"context"
	"fmt"
	"testing"
	"users/src"
	"users/src/domain"

	"github.com/gkampitakis/go-snaps/match"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func suite(t *testing.T) *repository {
	t.Helper()
	ctx := t.Context()

	pgContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)

	pgURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	repo, err := New(ctx, &src.Config{PostgresUrl: pgURL})
	require.NoError(t, err)

	return repo.(*repository)
}

func (repo *repository) suite(t *testing.T, users int) []*domain.User {
	t.Helper()
	ctx := t.Context()

	_, err := repo.pool.Exec(ctx, `DELETE FROM users`)
	require.NoError(t, err)

	var dummyUsers []*domain.User
	for i := range users {
		u := repoUser{
			ID:           uuid.New(),
			Username:     fmt.Sprintf("username%v", i),
			Email:        fmt.Sprintf("email%v", i),
			PasswordHash: fmt.Appendf([]byte{}, "passwordhash%v", i),
			PasswordSalt: fmt.Sprintf("passwordsalt%v", i),
		}

		_, err := repo.pool.Exec(ctx,
			`INSERT INTO users (id, username, email, password_hash, password_salt) VALUES ($1, $2, $3, $4, $5)`,
			u.ID, u.Username, u.Email, u.PasswordHash, u.PasswordSalt)
		require.NoError(t, err)

		dummyUsers = append(dummyUsers, u.domainify())
	}

	t.Cleanup(func() {
		rows, err := repo.pool.Query(context.Background(), `SELECT * FROM users`)
		require.NoError(t, err)
		usersResults, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[repoUser])
		require.NoError(t, err)
		snaps.MatchJSON(t, usersResults, match.Any("#.ID", "#.CreatedAt", "#.UpdatedAt"))
	})

	return dummyUsers
}
