package cmd

import (
	"errors"
	"io"

	"gabe565.com/changelog-generator/internal/config"
	"gabe565.com/changelog-generator/internal/git"
	"github.com/spf13/cobra"
)

func New(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "changelog-generator",
		Short:   "Generates a changelog from commits since the previous release",
		RunE:    run,
		Version: buildVersion(version, commit),

		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		DisableAutoGenTag: true,
	}
	config.RegisterFlags(cmd)
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	completionFlag, err := cmd.Flags().GetString(config.CompletionFlag)
	if err != nil {
		return err
	}
	if completionFlag != "" {
		return completion(cmd, args)
	}

	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	repo, err := git.FindRepo(cmd)
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
