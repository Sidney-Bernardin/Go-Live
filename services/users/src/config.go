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
	HttpSessionCookieDomain string `split_words:"true"`

	SessionDuration             time.Duration `split_words:"true" default:"24h"`
	ProfilePictureMaxBytes      int           `split_words:"true" default:"512000"`
	ProfilePictureInsertTimeout time.Duration `split_words:"true" default:"30s"`

	PostgresUrl     string `split_words:"true" required:"true"`
	RedisAddr       string `split_words:"true" required:"true"`
	AWSBaseRegion   string `split_words:"true" default:"us-east-1"`
	AWSBaseEndpoint string `split_words:"true"`
}

func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("app", &config)
	return &config, errors.Wrap(err, "cannot process config")
}
