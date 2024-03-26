package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, NewDefault())
}
