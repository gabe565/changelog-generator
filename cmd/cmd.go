package cmd

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/gabe565/changelog-generator/internal/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "changelog-generator",
		RunE: run,

		DisableAutoGenTag: true,
	}

	cmd.PersistentFlags().String("config", "", `Config file (default ".changelog-generator.yaml")`)
	_ = cmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	})

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return err
	}

	tags, err := repo.Tags()
	if err != nil {
		return err
	}

	var latest *plumbing.Reference
	var previous *plumbing.Reference
	if err := tags.ForEach(func(reference *plumbing.Reference) error {
		previous = latest
		latest = reference
		return nil
	}); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	tags.Close()

	head, err := repo.Reference(plumbing.HEAD, true)
	if err != nil {
		return err
	}
	if head.Hash() != latest.Hash() {
		previous = latest
		latest = head
	}

	commits, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return err
	}

	for {
		ref, err := commits.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		if previous != nil && ref.Hash == previous.Hash() {
			break
		}

		if !conf.Filters.Match(ref) {
			continue
		}

		for _, g := range conf.Groups {
			if g.Matches(ref) {
				g.AddCommit(ref)
				break
			}
		}
	}
	commits.Close()

	fmt.Println("## Changelog")
	slices.SortStableFunc(conf.Groups, func(a, b *config.Group) int {
		return a.Order - b.Order
	})
	var hasPrinted bool
	for _, g := range conf.Groups {
		g.Sort()
		if s := g.String(); s != "" {
			if hasPrinted && conf.Divider != "" {
				fmt.Println(conf.Divider)
			}
			hasPrinted = true
			fmt.Print(s)
		}
	}
	return nil
}
