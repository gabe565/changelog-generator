package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gabe565.com/changelog-generator/internal/config"
	"gabe565.com/changelog-generator/internal/git"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(opts ...cobrax.Option) *cobra.Command {
	name := "changelog-generator"
	var annotations map[string]string
	if base := filepath.Base(os.Args[0]); strings.HasPrefix(base, "git-") {
		// Installed as a git plugin
		name = base
		annotations = map[string]string{
			cobra.CommandDisplayNameAnnotation: strings.Replace(base, "-", " ", 1),
		}
	}

	cmd := &cobra.Command{
		Use:   name,
		Short: "Generates a changelog from commits since the previous release",
		RunE:  run,

		Annotations:       annotations,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		DisableAutoGenTag: true,
	}
	config.RegisterFlags(cmd)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	repo, err := git.PlainOpen(cmd)
	if err != nil {
		return err
	}

	previous, err := git.FindPreviousTag(repo, conf)
	if err != nil && !errors.Is(err, git.ErrNoPreviousTag) && !errors.Is(err, git.ErrNoCommits) {
		return err
	}

	if err := git.WalkCommits(repo, conf, previous); err != nil && !errors.Is(err, git.ErrNoCommits) {
		return err
	}

	_, err = io.WriteString(cmd.OutOrStdout(), conf.String())
	return err
}

func buildVersion(version, commit string) string {
	if commit != "" {
		version += " (" + commit + ")"
	}
	return version
}
