package cache

import (
	"context"
	"users/src/config"
	"users/src/domain/service"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cache struct {
	cfg *config.Config

	client *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (service.CacheRepository, error) {

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &cache{cfg, client}, nil
}
