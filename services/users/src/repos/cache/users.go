package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cacheUser struct {
	ID           domain.UUID `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	PasswordHash []byte      `json:"password_hash"`
	PasswordSalt string      `json:"password_salt"`
}

func (u *cacheUser) domainify() *domain.User {
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

func (c *cache) InsertUser(ctx context.Context, u *domain.User) error {
	key := fmt.Sprintf("user:%s", u.ID)
	err := c.client.JSONSet(ctx, key, ".", &cacheUser{
		ID:           u.ID,
		Username:     string(u.Username),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		PasswordSalt: string(u.PasswordSalt),
	}).Err()

	return errors.WithStack(err)
}

func (c *cache) GetUser(ctx context.Context, userID domain.UUID) (*domain.User, error) {
	key := fmt.Sprintf("user:%s", userID)

	userJSON, err := c.client.JSONGet(ctx, key, ".").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, service.ErrUserNotFound
		}

		return nil, errors.Wrap(err, "cannot get")
	}

	var user cacheUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errors.Wrap(err, "cannot decode")
	}

	return user.domainify(), nil
}
