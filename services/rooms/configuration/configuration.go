package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Configuration struct {
	HTTPPort        int           `required:"true" split_words:"true"`
	HTTPPongTimeout time.Duration `required:"true" split_words:"true"`

	UsersGRPCUrl string `required:"true" split_words:"true"`
	RedisURL     string `required:"true" split_words:"true"`
	RedisPassw   string `required:"true" split_words:"true"`
}

// New returns a pointer to a Configuration populated by environment variables.
func New(prefix string) (*Configuration, error) {

	var c Configuration

	// Populate the configuration with environment variables.
	if err := envconfig.Process(prefix, &c); err != nil {
		return nil, errors.Wrap(err, "cannot process configuration")
	}

	return &c, nil
}
