package token

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	CommandTokenArgsNumberMin = 2
	CommandTokenArgsNumberMax = 3
)

func validateTokenArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(CommandTokenArgsNumberMin)(cmd, args); err != nil {
		return errors.Wrap(errValidation, err.Error())
	}

	if err := cobra.MaximumNArgs(CommandTokenArgsNumberMax)(cmd, args); err != nil {
		return errors.Wrap(errValidation, err.Error())
	}

	return nil
}
