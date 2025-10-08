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

//go:embed Migrations/*.sql
var migrations embed.FS

type databaseRepository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, config *src.Config) (service.DatabaseRepository, error) {

	pool, err := pgxpool.New(ctx, config.PostgresUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating connection pool")
	}

	if err := doMigrations(config); err != nil {
		return nil, errors.Wrap(err, "failed migrations")
	}

	return &databaseRepository{pool}, nil
}

func doMigrations(config *src.Config) error {

	source, err := iofs.New(migrations, "Migrations")
	if err != nil {
		return errors.Wrap(err, "failed creating migration source")
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, config.PostgresUrl)
	if err != nil {
		return errors.Wrap(err, "failed creating migrator")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed migrating")
	}

	return nil
}
