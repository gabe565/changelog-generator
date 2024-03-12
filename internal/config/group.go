package config

import (
	"regexp"
	"slices"
	"strings"

	"github.com/gabe565/changelog-generator/internal/util"
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
		return g.re.MatchString(c.Message)
	}
	return true
}

func (g *Group) AddCommit(c *object.Commit) {
	g.Commits = append(g.Commits, c)
}

func (g *Group) Sort() {
	switch Default.Sort {
	case SortAscending:
		slices.SortStableFunc(g.Commits, func(a, b *object.Commit) int {
			return strings.Compare(a.Message, b.Message)
		})
	case SortDescending:
		slices.SortStableFunc(g.Commits, func(a, b *object.Commit) int {
			return strings.Compare(b.Message, a.Message)
		})
	}
}

func (g *Group) String() string {
	if len(g.Commits) == 0 {
		return ""
	}

	var result string
	if g.Title != "" {
		result += "### " + g.Title + "\n"
	}

	for _, commit := range g.Commits {
		entry := "- "
		if Default.Abbrev != HashDisabled {
			entry += commit.Hash.String()[:Default.Abbrev] + " "
		}
		entry += util.ShortMessage(commit)
		result += entry + "\n"
	}
	return result
}
