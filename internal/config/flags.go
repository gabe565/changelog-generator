package config

import (
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const (
	ConfigFlag     = "config"
	RepoFlag       = "repo"
	CompletionFlag = "completion"
)

func RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().String(ConfigFlag, "", `Config file (default ".changelog-generator.yaml")`)
	must.Must(cmd.RegisterFlagCompletionFunc(ConfigFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	}))

	cmd.Flags().String(RepoFlag, ".", `Path to the git repo root. Parent directories will be walked until .git is found.`)
	must.Must(cmd.RegisterFlagCompletionFunc(RepoFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveFilterDirs
	}))

	cmd.Flags().String(CompletionFlag, "", "Output command-line completion code for the specified shell. (one of bash, zsh, fish, powershell)")
	must.Must(cmd.RegisterFlagCompletionFunc(CompletionFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
	}))
}
