package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	SessionDuration time.Duration `split_words:"true" required:"true"`

	PostgresUrl string `split_words:"true" required:"true"`
	RedisAddr   string `split_words:"true" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("app", &cfg)
	return &cfg, errors.Wrap(err, "cannot process config")
}
