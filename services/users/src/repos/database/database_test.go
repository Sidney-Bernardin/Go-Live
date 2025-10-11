package database

import (
	"testing"
	"users/src"
	"users/src/domain"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newTestDatabase(t *testing.T) *databaseRepository {
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

	db, err := New(ctx, &src.Config{PostgresUrl: connStr})
	require.NoErrorf(t, err, "Failed creating repository: %v", err)

	return db.(*databaseRepository)
}

func newTestUser(t *testing.T) *domain.User {
	t.Helper()

	ctx := t.Context()
	u := domain.User{}

	u.ID = domain.NewUUID()
	u.Email = gofakeit.Email()
	u.PasswordSalt = domain.NewPasswordSalt()

	password, err := domain.NewPassword(ctx, domain.MustRandomString(16))
	require.NoError(t, err)

	u.Username, err = domain.NewUsername(ctx, gofakeit.Username())
	require.NoError(t, err)

	u.PasswordHash, err = domain.NewPasswordHash(ctx, password, u.PasswordSalt)
	require.NoError(t, err)

	return &u
}

func (db *databaseRepository) assertUser(t *testing.T, userID domain.UUID, user *domain.User) {
	t.Helper()
	ctx := t.Context()

	rows, err := db.pool.Query(ctx, `SELECT id, username, email, password_hash, password_salt FROM users WHERE id = $1`, userID)
	require.NoError(t, err)
	dbuser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dbUser])
	u := dbuser.domainify()

	if user == nil {
		require.Equal(t, pgx.ErrNoRows, err)
		return
	}

	require.NoError(t, err)
	assert.Equal(t, user.ID, u.ID)
	assert.Equal(t, user.Username, u.Username)
	assert.Equal(t, user.Email, u.Email)
	assert.Equal(t, user.PasswordHash, u.PasswordHash)
	assert.Equal(t, user.PasswordSalt, u.PasswordSalt)
}
