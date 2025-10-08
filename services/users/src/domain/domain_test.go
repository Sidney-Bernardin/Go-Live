package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	t.Parallel()

	t.Run("Length", func(t *testing.T) {
		assert.Equal(t, len(MustRandomString(32)), 32)
		assert.Equal(t, len(MustRandomString(33)), 32)
	})

	t.Run("Randomness", func(t *testing.T) {
		for range 1000000 {
			require.NotEqual(t, MustRandomString(32), MustRandomString(32))
		}
	})
}
