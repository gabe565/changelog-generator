package cmd

import (
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

	registerCompletionFlag(cmd)

	cmd.Flags().String("config", "", `Config file (default ".changelog-generator.yaml")`)
	_ = cmd.RegisterFlagCompletionFunc("config", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	})

	cmd.Flags().String("repo", ".", `Path to the git repo root. Parent directories will be walked until .git is found.`)
	_ = cmd.RegisterFlagCompletionFunc("repo", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveFilterDirs
	})

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	completionFlag, err := cmd.Flags().GetString(CompletionFlag)
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
	if err != nil {
		return err
	}

	if err := git.WalkCommits(repo, conf, previous); err != nil {
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
