package src_test

import (
	"database/sql"
	"testing"
	"users/src"
	"users/src/domain"
	"users/src/domain/service"
	"users/src/repos/database"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newCacheRepository(t *testing.T) (service.CacheRepository, *redis.Client) {
	return nil, nil
}
