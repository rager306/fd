package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortHash(t *testing.T) {
	assert.Equal(t, "b94d27b9934d", shortHash("hello world"))
	assert.Equal(t, "e3b0c44298fc", shortHash(""))
}
