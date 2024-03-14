package util

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestShortMessage(t *testing.T) {
	type args struct {
		c *object.Commit
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"already short", args{&object.Commit{Message: "test"}}, "test"},
		{"long", args{&object.Commit{Message: "test\n\ntest"}}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ShortMessage(tt.args.c))
		})
	}
}
