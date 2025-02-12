package config

import (
	"cmp"
	"regexp"
	"slices"

	"gabe565.com/changelog-generator/internal/util"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	SortAscending  = "asc"
	SortDescending = "desc"

	HashDisabled = -1
)

type Group struct {
	Title  string `yaml:"title"`
	Order  int    `yaml:"order"`
	Regexp string `yaml:"regexp"`
	re     *regexp.Regexp

	Commits []*object.Commit `yaml:"-"`
}

func (g *Group) Matches(c *object.Commit) bool {
	if g.re != nil {
		shortMessage := util.ShortMessage(c)
		return g.re.MatchString(shortMessage)
	}
	return true
}

func (g *Group) AddCommit(c *object.Commit) {
	g.Commits = append(g.Commits, c)
}

func (g *Group) Sort(sort string) {
	switch sort {
	case SortAscending:
		slices.SortStableFunc(g.Commits, func(a, b *object.Commit) int {
			return cmp.Compare(a.Message, b.Message)
		})
	case SortDescending:
		slices.SortStableFunc(g.Commits, func(a, b *object.Commit) int {
			return cmp.Compare(b.Message, a.Message)
		})
	}
}

func (g *Group) String(conf *Config) string {
	if len(g.Commits) == 0 {
		return ""
	}

	var result string
	if g.Title != "" {
		result += "### " + g.Title + "\n"
	}

	for _, commit := range g.Commits {
		entry := "- "
		if conf.Abbrev != HashDisabled {
			entry += commit.Hash.String()[:conf.Abbrev] + " "
		}
		entry += util.ShortMessage(commit)
		result += entry + "\n"
	}
	return result
}
