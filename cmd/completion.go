package cmd

import (
	"errors"
	"fmt"

	"github.com/gabe565/changelog-generator/internal/config"
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
		if err := cmd.Root().GenBashCompletion(cmd.OutOrStdout()); err != nil {
			return err
		}
	case "zsh":
		if err := cmd.Root().GenZshCompletion(cmd.OutOrStdout()); err != nil {
			return err
		}
	case "fish":
		if err := cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true); err != nil {
			return err
		}
	case "powershell":
		if err := cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidShell, completionFlag)
	}
	return nil
}
