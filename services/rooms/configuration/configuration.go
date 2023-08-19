package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Port            int           `required:"true" split_words:"true"`
	ReadTimeout     time.Duration `required:"true" split_words:"true"`
	WriteTimeout    time.Duration `required:"true" split_words:"true"`
	WSCloseTimeout  time.Duration `required:"true" split_words:"true"`
	ShutdownTimeout time.Duration `required:"true" split_words:"true"`

	UsersGRPCUrl string `required:"true" split_words:"true"`
	CacheURL     string `required:"true" split_words:"true"`
	CachePassw   string `required:"true" split_words:"true"`
}

// NewConfig returns a Config populated with environment variables.
func NewConfig(prefix string) (*Config, error) {
	var c Config
	err := envconfig.Process(prefix, &c)
	return &c, errors.Wrap(err, "cannot generate configuration")
}
