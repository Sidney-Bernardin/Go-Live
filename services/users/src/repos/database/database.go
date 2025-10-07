package database

import (
	"context"
	"embed"
	"users/src/config"
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

func New(ctx context.Context, cfg *config.Config) (service.DatabaseRepository, error) {

	pool, err := pgxpool.New(ctx, cfg.PostgresUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating connection pool")
	}

	if err := doMigrations(cfg); err != nil {
		return nil, errors.Wrap(err, "failed migrations")
	}

	return &databaseRepository{pool}, nil
}

func doMigrations(cfg *config.Config) error {

	source, err := iofs.New(migrations, "Migrations")
	if err != nil {
		return errors.Wrap(err, "failed creating migration source")
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, cfg.PostgresUrl)
	if err != nil {
		return errors.Wrap(err, "failed creating migrator")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed migrating")
	}

	return nil
}

type rowModel[D any] interface {
	domainify() *D
}

func domainifyMany[R rowModel[D], D any](rr []R) []*D {
	dd := make([]*D, len(rr))
	for i, r := range rr {
		dd[i] = r.domainify()
	}
	return dd
}
