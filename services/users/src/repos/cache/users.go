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

type repoUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *repoUser) domainify() *domain.User {
	if u == nil {
		return nil
	}

	return &domain.User{
		ID:        domain.UUID(u.ID),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  domain.Username(u.Username),
		Email:     u.Email,
	}
}

func (repo *repository) InsertUser(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user:%s", user.ID)
	err := repo.client.JSONSet(ctx, key, ".", &repoUser{
		ID:        uuid.UUID(user.ID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  string(user.Username),
		Email:     user.Email,
	}).Err()

	return errors.Wrap(err, "failed setting")
}

func (repo *repository) GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", userID)

	// Get the user.
	userJSON, err := repo.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "failed getting")
	}

	// Decode the user.
	var user repoUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errors.Wrap(err, "failed decoding")
	}

	return user.domainify(), nil
}
