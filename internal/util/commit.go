package util

import (
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func ShortMessage(c *object.Commit) string {
	shortMessage, _, _ := strings.Cut(c.Message, "\n")
	return shortMessage
}
