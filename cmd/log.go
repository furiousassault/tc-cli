package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

var cmdLog = &cobra.Command{
	Use:   "log <build id>",
	Short: "Show logs for of the build specified by id",
	Args:  validateLogArgs,
	RunE:   buildLog,
}

func validateLogArgs(cmd *cobra.Command, args []string) error {

	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
		return fmt.Errorf("enter valid positive integer buildId")
	}

	return nil
}

func buildLog(_ *cobra.Command, args []string) error {
	out, err := api.API().Logs.GetBuildLog(args[0])
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
