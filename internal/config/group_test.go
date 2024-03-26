package config

import (
	"regexp"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestGroup_Matches(t *testing.T) {
	t.Parallel()
	type fields struct {
		Title   string
		Order   int
		Regexp  string
		re      *regexp.Regexp
		Commits []*object.Commit
	}
	type args struct {
		c *object.Commit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"null regexp", fields{}, args{&object.Commit{Message: "test"}}, true},
		{"has regexp match", fields{re: regexp.MustCompile("test")}, args{&object.Commit{Message: "test"}}, true},
		{"no regexp match", fields{re: regexp.MustCompile("example")}, args{&object.Commit{Message: "test"}}, false},
		{"only match first line", fields{re: regexp.MustCompile("test")}, args{&object.Commit{Message: "example\ntest"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			g := &Group{
				Title:   tt.fields.Title,
				Order:   tt.fields.Order,
				Regexp:  tt.fields.Regexp,
				re:      tt.fields.re,
				Commits: tt.fields.Commits,
			}
			assert.Equal(t, tt.want, g.Matches(tt.args.c))
		})
	}
}

func TestGroup_String(t *testing.T) {
	t.Parallel()
	testCommit := &object.Commit{Message: "test", Hash: plumbing.NewHash("DEADBEEF")}

	type fields struct {
		Title   string
		Order   int
		Regexp  string
		re      *regexp.Regexp
		Commits []*object.Commit
	}
	tests := []struct {
		name   string
		abbrev int
		fields fields
		want   string
	}{
		{"no commits", 8, fields{Title: "Test"}, ""},
		{"no title", 8, fields{Commits: []*object.Commit{testCommit}}, "- deadbeef test\n"},
		{"title and commits", 8, fields{Title: "Test", Commits: []*object.Commit{testCommit}}, "### Test\n- deadbeef test\n"},
		{"skip commit hash", -1, fields{Commits: []*object.Commit{testCommit}}, "- test\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewDefault()
			c.Abbrev = tt.abbrev

			g := &Group{
				Title:   tt.fields.Title,
				Order:   tt.fields.Order,
				Regexp:  tt.fields.Regexp,
				re:      tt.fields.re,
				Commits: tt.fields.Commits,
			}
			assert.Equal(t, tt.want, g.String(c))
		})
	}
}
