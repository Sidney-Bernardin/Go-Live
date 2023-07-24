package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Configuration struct {
	HTTPPort int `required:"true" split_words:"true"`
	GRPCPort int `required:"true" split_words:"true"`

	MongoURL string `required:"true" split_words:"true"`

	DatabaseTimeout time.Duration `required:"true" split_words:"true"`
	SessionLength   time.Duration `required:"true" split_words:"true"`
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
