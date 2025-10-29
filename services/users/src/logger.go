package src

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/pkg/errors"
)

func NewLogger(config *Config) *slog.Logger {
	var h slog.Handler

	switch {
	case config.LogPretty:
		h = devslog.NewHandler(os.Stderr, &devslog.Options{
			NoColor: !config.LogPrettyColor,
		})

	default:
		h = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	}

	return slog.New(h)
}

func ErrGroup(e error) slog.Attr {
	return slog.Group("err",
		"type", fmt.Sprintf("%T", errors.Cause(e)),
		"msg", e.Error(),
		"stack", fmt.Sprintf("%+v", e),
	)
}
