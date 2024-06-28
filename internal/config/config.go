package config

import (
	"slices"
)

type Config struct {
	Tag     Tag      `yaml:"tag"`
	Sort    string   `yaml:"sort"`
	Abbrev  int      `yaml:"abbrev"`
	Groups  []*Group `yaml:"groups"`
	Divider string   `yaml:"divider"`
	Filters Filters  `yaml:"filters"`
}

func NewDefault() *Config {
	return &Config{
		Sort:   SortAscending,
		Abbrev: 8,
		Groups: []*Group{{}},
	}
}

func (c *Config) String() string {
	var result string
	result += "## Changelog\n"
	slices.SortStableFunc(c.Groups, func(a, b *Group) int {
		return a.Order - b.Order
	})
	var hasPrinted bool
	for _, g := range c.Groups {
		g.Sort(c.Sort)
		if s := g.String(c); s != "" {
			if hasPrinted && c.Divider != "" {
				result += c.Divider + "\n"
			}
			hasPrinted = true
			result += s
		}
	}
	return result
}
