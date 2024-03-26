package config

import (
	"regexp"
	"testing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestFilters_Match(t *testing.T) {
	t.Parallel()
	type fields struct {
		Exclude   []string
		excludeRe []*regexp.Regexp
		Include   []string
		includeRe []*regexp.Regexp
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
		{"no filters", fields{}, args{&object.Commit{Message: "test"}}, true},
		{
			"include filter match",
			fields{includeRe: []*regexp.Regexp{regexp.MustCompile("test")}},
			args{&object.Commit{Message: "test"}},
			true,
		},
		{
			"include filter no match",
			fields{includeRe: []*regexp.Regexp{regexp.MustCompile("example")}},
			args{&object.Commit{Message: "test"}},
			false,
		},
		{
			"exclude filter match",
			fields{excludeRe: []*regexp.Regexp{regexp.MustCompile("test")}},
			args{&object.Commit{Message: "test"}},
			false,
		},
		{
			"exclude filter no match",
			fields{excludeRe: []*regexp.Regexp{regexp.MustCompile("example")}},
			args{&object.Commit{Message: "test"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &Filters{
				Exclude:   tt.fields.Exclude,
				excludeRe: tt.fields.excludeRe,
				Include:   tt.fields.Include,
				includeRe: tt.fields.includeRe,
			}
			assert.Equal(t, tt.want, f.Match(tt.args.c))
		})
	}
}
