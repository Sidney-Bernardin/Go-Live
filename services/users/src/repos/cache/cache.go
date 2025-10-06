package cache

import (
	"context"
	"encoding/json"
	"users/src/config"
	"users/src/domain"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cache struct {
	cfg *config.Config

	client *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (domain.CacheRepository, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &cache{cfg, client}, nil
}

func jsonGet[M any](ctx context.Context, client redis.JSONCmdable, key string, paths ...string) (*M, error) {

	modelJSON, err := client.JSONGet(ctx, key, paths...).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get")
	}

	var model M
	if err := json.Unmarshal([]byte(modelJSON), &model); err != nil {
		return nil, errors.Wrap(err, "cannot decode")
	}

	return &model, nil
}
