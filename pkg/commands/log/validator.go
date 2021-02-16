package log

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func validateLogArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
		return fmt.Errorf("%w: enter valid positive integer buildId", err)
	}

	return nil
}
