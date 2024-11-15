package config

import (
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	FlagConfig     = "config"
	FlagRepo       = "repo"
	FlagCompletion = "completion"
)

func RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagConfig, "", `Config file (default ".changelog-generator.yaml")`)
	must.Must(cmd.RegisterFlagCompletionFunc(FlagConfig, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	}))

	cmd.Flags().String(FlagRepo, ".", `Path to the git repo root. Parent directories will be walked until .git is found.`)
	must.Must(cmd.RegisterFlagCompletionFunc(FlagRepo, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveFilterDirs
	}))

	cmd.Flags().String(FlagCompletion, "", "Output command-line completion code for the specified shell. (one of bash, zsh, fish, powershell)")
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCompletion, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
	}))
}
