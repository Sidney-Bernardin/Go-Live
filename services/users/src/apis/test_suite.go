package apis

import (
	"bytes"
	"log/slog"
	"testing"
	"users/src"
	"users/src/domain/service"
	"users/src/repos/cache"
	"users/src/repos/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type APITestSuite struct {
	LogBuf *bytes.Buffer

	Config *src.Config
	Logger *slog.Logger

	DatabaseClient *sqlx.DB
	CacheClient    *redis.Client

	Service *service.Service
}

func NewAPITestSuite(t *testing.T) *APITestSuite {
	t.Helper()

	ctx := t.Context()
	var err error

	suite := &APITestSuite{}

	suite.Config, err = src.NewConfig()
	require.NoError(t, err)

	suite.LogBuf = bytes.NewBuffer(make([]byte, 0))
	suite.Logger = slog.New(slog.NewJSONHandler(suite.LogBuf, nil))

	databaseRepo, err := database.New(ctx, suite.Config)
	require.NoError(t, err)

	suite.DatabaseClient, err = sqlx.Open("pgx", suite.Config.PostgresUrl)
	require.NoError(t, err)

	cacheRepo, err := cache.New(ctx, suite.Config)
	require.NoError(t, err)

	suite.CacheClient = redis.NewClient(&redis.Options{Addr: suite.Config.RedisAddr})
	err = suite.CacheClient.Ping(ctx).Err()
	require.NoError(t, err)

	suite.Service = service.New(suite.Config, databaseRepo, cacheRepo)
	return suite
}
