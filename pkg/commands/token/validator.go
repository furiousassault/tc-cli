package token

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const CommandTokenArgsNumber = 3

func validateTokenArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(CommandTokenArgsNumber)(cmd, args); err != nil {
		return errors.Wrap(errValidation, err.Error())
	}

	if args[1] == args[2] {
		return errors.Wrap(errValidation, "old and new token names must not be equal")
	}

	return nil
}
