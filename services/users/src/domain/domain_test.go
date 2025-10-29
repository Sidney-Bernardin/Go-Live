package domain

import (
	"testing"
	"users/src"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	t.Parallel()

	t.Run("length", func(t *testing.T) {
		assert.Equal(t, len(src.MustRandomString(32)), 32)
		assert.Equal(t, len(src.MustRandomString(33)), 32)
	})

	t.Run("randomness", func(t *testing.T) {
		for range 1000000 {
			require.NotEqual(t, src.MustRandomString(32), src.MustRandomString(32))
		}
	})
}
