package test

import (
	"database/sql"
	"testing"
	"users/src"
	"users/src/domain/service"
	"users/src/repos/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewDatabaseRepository(t *testing.T) (service.DatabaseRepository, *sql.DB) {
	t.Helper()
	ctx := t.Context()

	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies())

	t.Cleanup(func() {
		err := testcontainers.TerminateContainer(container)
		assert.NoErrorf(t, err, "Failed terminating container: %v", err)
	})

	require.NoErrorf(t, err, "Failed running container: %v", err)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoErrorf(t, err, "Failed getting container connection string: %v", err)

	repo, err := database.New(ctx, &src.Config{PostgresUrl: connStr})
	require.NoErrorf(t, err, "Failed creating repository: %v", err)

	conn, err := sql.Open("", "")
	require.NoErrorf(t, err, "Failed opening database: %v", err)

	return repo, conn
}
