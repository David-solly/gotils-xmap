package xmap

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestXmap(t *testing.T) {
	t.Run("Test gotils xmap", func(t *testing.T) {
		assert.Equal(t, hi(), true)
	})
}

func hi() bool {
	return true
}
