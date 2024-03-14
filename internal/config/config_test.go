package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefault(t *testing.T) {
	assert.NotNil(t, NewDefault())
}
