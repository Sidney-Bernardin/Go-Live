package cache

import (
	"context"
	"users/src"
	"users/src/domain/service"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type cache struct {
	config *src.Config

	client *redis.Client
}

func New(ctx context.Context, config *src.Config) (service.CacheRepository, error) {

	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "ping failed")
	}

	return &cache{config, client}, nil
}
