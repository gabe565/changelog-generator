package config

import "github.com/spf13/cobra"

const (
	ConfigFlag     = "config"
	RepoFlag       = "repo"
	CompletionFlag = "completion"
)

func RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().String(ConfigFlag, "", `Config file (default ".changelog-generator.yaml")`)
	if err := cmd.RegisterFlagCompletionFunc(ConfigFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml"}, cobra.ShellCompDirectiveFilterFileExt
	}); err != nil {
		panic(err)
	}

	cmd.Flags().String(RepoFlag, ".", `Path to the git repo root. Parent directories will be walked until .git is found.`)
	if err := cmd.RegisterFlagCompletionFunc(RepoFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{}, cobra.ShellCompDirectiveFilterDirs
	}); err != nil {
		panic(err)
	}

	cmd.Flags().String(CompletionFlag, "", "Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.")
	if err := cmd.RegisterFlagCompletionFunc(CompletionFlag, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		panic(err)
	}
}
