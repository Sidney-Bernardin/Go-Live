package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port            int           `default:"8080" split_words:"true"`
	ReadTimeout     time.Duration `default:"15s" split_words:"true"`
	WriteTimeout    time.Duration `default:"15s" split_words:"true"`
	WSCloseTimeout  time.Duration `default:"5s" split_words:"true"`
	ShutdownTimeout time.Duration `default:"5s" split_words:"true"`

	UsersGRPCUrl string `default:"localhost:9000" split_words:"true"`
	CacheURL     string `default:"localhost:6379" split_words:"true"`
	CachePassw   string `default:"" split_words:"true"`
}

// NewConfig returns a Config populated with environment variables.
func NewConfig(prefix string) (*Config, error) {
	var c Config
	err := envconfig.Process(prefix, &c)
	return &c, errors.Wrap(err, "cannot generate configuration")
}
