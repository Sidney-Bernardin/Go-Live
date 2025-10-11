package database

import (
	"context"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type dbUser struct {
	ID domain.UUID `db:"id"`

	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	PasswordSalt string `db:"password_salt"`
}

func (u *dbUser) domainify() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{
		ID:           u.ID,
		Username:     domain.Username(u.Username),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		PasswordSalt: domain.PasswordSalt(u.PasswordSalt),
	}
}

const qInsertUser = `
	INSERT INTO users (id, username, email, password_hash, password_salt)
	VALUES ($1, $2, $3, $4, $5)
`

func (repo *databaseRepository) InsertUser(ctx context.Context, user *domain.User) error {
	_, err := repo.pool.Exec(ctx, qInsertUser,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.PasswordSalt)

	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):

			switch pgErr.ConstraintName {
			case "users_username_key":
				switch pgErr.Code {
				case "23505":
					return service.ErrUsernameTaken
				}

			case "users_email_key":
				switch pgErr.Code {
				case "23505":
					return service.ErrEmailTaken
				}
			}
		}

		return errors.WithStack(err)
	}

	return nil
}

const qGetUserByID = `
	SELECT * FROM users WHERE id = $1
`

func (repo *databaseRepository) GetUserByID(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	rows, err := repo.pool.Query(ctx, qGetUserByID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed selecting")
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dbUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "failed collecting rows")
	}

	return user.domainify(), nil
}
