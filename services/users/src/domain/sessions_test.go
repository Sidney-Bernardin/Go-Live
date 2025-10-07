package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSRFToken(t *testing.T) {
	t.Parallel()

	t.Run("Length", func(t *testing.T) {
		assert.Equal(t, len(NewPasswordSalt()), 32)
	})
}
