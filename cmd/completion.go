package cmd

import (
	"errors"
	"fmt"

	"gabe565.com/changelog-generator/internal/config"
	"github.com/spf13/cobra"
)

var ErrInvalidShell = errors.New("invalid shell")

func completion(cmd *cobra.Command, _ []string) error {
	completionFlag, err := cmd.Flags().GetString(config.CompletionFlag)
	if err != nil {
		panic(err)
	}

	switch completionFlag {
	case "bash":
		return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
	case "zsh":
		return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
	case "fish":
		return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
	default:
		return fmt.Errorf("%w: %s", ErrInvalidShell, completionFlag)
	}
}
