package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type userRecord struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

func (r *userRecord) user() *domain.User {
	if r == nil {
		return nil
	}

	return &domain.User{
		ID:        domain.UUID(r.ID),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Username:  domain.Username(r.Username),
		Email:     r.Email,
	}
}

func userKey(userID domain.UUID) string {
	return fmt.Sprintf("user:%s", userID)
}

func (repo *repository) InsertUser(ctx context.Context, user *domain.User) error {
	err := repo.client.JSONSet(ctx, userKey(user.ID), ".", &userRecord{
		ID:        uuid.UUID(user.ID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  string(user.Username),
		Email:     user.Email,
	}).Err()

	return errors.Wrap(err, "failed setting")
}

func (repo *repository) GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error) {

	// Get the user.
	userRecordJSON, err := repo.client.JSONGet(ctx, userKey(userID), ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "failed getting")
	}

	// Decode the user.
	var userRecord userRecord
	if err := json.Unmarshal([]byte(userRecordJSON), &userRecord); err != nil {
		return nil, errors.Wrap(err, "failed decoding")
	}

	return userRecord.user(), nil
}
