package src

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	LogPretty      bool `split_words:"true" default:"true"`
	LogPrettyColor bool `split_words:"true" default:"true"`

	HttpAddr                string `split_words:"true" required:"true"`
	HttpSessionCookieDomain string `split_words:"true" required:"true"`

	SessionDuration time.Duration `split_words:"true" required:"true"`

	PostgresUrl string `split_words:"true" required:"true"`
	RedisAddr   string `split_words:"true" required:"true"`
}

func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("app", &config)
	return &config, errors.Wrap(err, "cannot process config")
}
