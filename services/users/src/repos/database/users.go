package database

import (
	"context"
	"time"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type userRow struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
	PasswordSalt string `db:"password_salt"`
}

func (r *userRow) user() *domain.User {
	if r == nil {
		return nil
	}

	return &domain.User{
		ID:           domain.UUID(r.ID),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		Username:     domain.Username(r.Username),
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		PasswordSalt: domain.PasswordSalt(r.PasswordSalt),
	}
}

const qInsertUser = `
	INSERT INTO users (id, username, email, password_hash, password_salt)
	VALUES ($1, $2, $3, $4, $5)
`

func (repo *repository) InsertUser(ctx context.Context, user *domain.User) error {
	_, err := repo.pool.Exec(ctx, qInsertUser,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.PasswordSalt)

	if pgErr := (&pgconn.PgError{}); errors.As(err, &pgErr) {
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

	return errors.Wrap(err, "failed inserting")
}

const qGetUserByID = `
	SELECT * FROM users WHERE id = $1
`

func (repo *repository) GetUserByID(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	// Get the user with a matching user-ID.
	userRows, err := repo.pool.Query(ctx, qGetUserByID, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed selecting")
	}

	// Decode the user.
	userRow, err := pgx.CollectExactlyOneRow(userRows, pgx.RowToAddrOfStructByName[userRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "failed collecting rows")
	}

	return userRow.user(), nil
}
