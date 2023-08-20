package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Mock            bool          `default:"false" split_words:"true"`
	ShutdownTimeout time.Duration `default:"5s" split_words:"true"`

	HTTPPort         int           `required:"true" split_words:"true"`
	HTTPReadTimeout  time.Duration `default:"15s" split_words:"true"`
	HTTPWriteTimeout time.Duration `default:"15s" split_words:"true"`
	GRPCPort         int           `required:"true" split_words:"true"`

	SessionLength time.Duration `default:"720m" split_words:"true"`

	DBConnTimeout time.Duration `default:"10s" split_words:"true"`
	MongoURL      string        `default:"mongodb://localhost:27017/" split_words:"true"`
}

// NewConfig returns a Config populated with environment variables.
func NewConfig(prefix string) (*Config, error) {
	var c Config
	err := envconfig.Process(prefix, &c)
	return &c, errors.Wrap(err, "cannot generate configuration")
}
