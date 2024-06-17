package config

import "regexp"

type Tag struct {
	Regexp string `yaml:"regexp"`
	re     *regexp.Regexp
}

func (t *Tag) Match(s string) bool {
	if t.re == nil {
		return true
	}
	return t.re.MatchString(s)
}
