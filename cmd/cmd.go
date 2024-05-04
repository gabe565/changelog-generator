package cmd

import (
	"errors"
	"io"

	"github.com/gabe565/changelog-generator/internal/config"
	"github.com/gabe565/changelog-generator/internal/git"
	"github.com/spf13/cobra"
)

func New(version, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "changelog-generator",
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

	previous, err := git.FindRefs(repo)
	if err != nil && !errors.Is(err, git.ErrNoPreviousTag) && !errors.Is(err, git.ErrNoCommits) {
		return err
	}

	if err := git.WalkCommits(repo, conf, previous); err != nil && !errors.Is(err, git.ErrNoCommits) {
		return err
	}

	_, _ = io.WriteString(cmd.OutOrStdout(), conf.String())
	return nil
}

func buildVersion(version, commit string) string {
	if commit != "" {
		version += " (" + commit + ")"
	}
	return version
}
