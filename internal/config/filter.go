package config

import (
	"regexp"

	"gabe565.com/changelog-generator/internal/util"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Filters struct {
	Exclude   []string `yaml:"exclude"`
	excludeRe []*regexp.Regexp
	Include   []string `yaml:"include"`
	includeRe []*regexp.Regexp
}

func (f *Filters) Match(c *object.Commit) bool {
	shortMessage := util.ShortMessage(c)
	switch {
	case len(f.includeRe) != 0:
		for _, re := range f.includeRe {
			if re.MatchString(shortMessage) {
				return true
			}
		}
		return false
	case len(f.excludeRe) != 0:
		for _, re := range f.excludeRe {
			if re.MatchString(shortMessage) {
				return false
			}
		}
		return true
	}
	return true
}
