package database

import (
	"context"
	"embed"
	"users/src"
	"users/src/domain/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

//go:embed Migrations
var migrationsDir embed.FS

type repository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, config *src.Config) (service.DatabaseRepository, error) {
	pool, err := pgxpool.New(ctx, config.PostgresUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating connection pool")
	}

	if err := migrations(config); err != nil {
		return nil, errors.Wrap(err, "failed migrations")
	}

	return &repository{pool}, nil
}

func migrations(config *src.Config) error {

	// Create migration source with the Migrations directory.
	source, err := iofs.New(migrationsDir, "Migrations")
	if err != nil {
		return errors.Wrap(err, "failed creating migration source")
	}

	// Create a migrate instance.
	m, err := migrate.NewWithSourceInstance("iofs", source, config.PostgresUrl)
	if err != nil {
		return errors.Wrap(err, "failed creating migrator")
	}

	// Migrate up.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed migrating")
	}

	return nil
}
