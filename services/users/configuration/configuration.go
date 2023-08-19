package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	ShutdownTimeout  time.Duration `required:"true" split_words:"true"`
	HTTPPort         int           `required:"true" split_words:"true"`
	HTTPReadTimeout  time.Duration `required:"true" split_words:"true"`
	HTTPWriteTimeout time.Duration `required:"true" split_words:"true"`
	GRPCPort         int           `required:"true" split_words:"true"`

	SessionLength time.Duration `required:"true" split_words:"true"`

	DatabaseTimeout time.Duration `required:"true" split_words:"true"`
	MongoURL        string        `required:"true" split_words:"true"`
}

// NewConfig returns a Config populated with environment variables.
func NewConfig(prefix string) (*Config, error) {
	var c Config
	err := envconfig.Process(prefix, &c)
	return &c, errors.Wrap(err, "cannot generate configuration")
}
