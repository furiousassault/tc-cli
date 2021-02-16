package token

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	CommandTokenArgsNumberMin = 2
	CommandTokenArgsNumberMax = 3
)

func validateTokenArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(CommandTokenArgsNumberMin)(cmd, args); err != nil {
		return fmt.Errorf("%w: %s", errValidation, err.Error())
	}

	if err := cobra.MaximumNArgs(CommandTokenArgsNumberMax)(cmd, args); err != nil {
		return fmt.Errorf("%w: %s", errValidation, err.Error())
	}

	return nil
}
