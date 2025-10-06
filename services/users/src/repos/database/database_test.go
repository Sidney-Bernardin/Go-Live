package database

import (
	"testing"
	"users/src/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newTest(t *testing.T) *databaseRepository {
	t.Helper()

	ctx := t.Context()

	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies())

	t.Cleanup(func() {
		err := testcontainers.TerminateContainer(container)
		assert.NoErrorf(t, err, "Failed terminating container: %v", err)
	})

	if err != nil {
		require.NoErrorf(t, err, "Failed running container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		require.NoErrorf(t, err, "Failed getting container connection string: %v", err)
	}

	repo, err := New(ctx, &config.Config{
		PostgresUrl: connStr,
	})

	if err != nil {
		require.NoErrorf(t, err, "Failed creating repository: %v", err)
	}

	return repo.(*databaseRepository)
}
