package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

var (
	cmdRun = &cobra.Command{
		Use:   "run <build_configuration_id>",
		Short: "run build configuration",
		Args:  cobra.ExactArgs(1),
		RunE:  buildConfigurationRun,
	}
)

func buildConfigurationRun(_ *cobra.Command, args []string) error {
	result, err := api.API().BuildQueue.RunBuildByBuildConfID(args[0])
	if err != nil {
		return err
	}

	output.Output().List(
		[]string{"BUILD_ID", "QUEUED_BY", "STATE"},
		[][]string{{fmt.Sprint(result.BuildID), result.TriggeredBy, result.BuildState}},
	)
	fmt.Println()
	fmt.Println("Properties")

	properties := make([][]string, 0)

	for _, property := range result.Parameters.Items {
		properties = append(properties, []string{property.Type, property.Name, property.Value, fmt.Sprint(property.Inherited)})
	}

	output.Output().List([]string{"TYPE", "KEY", "VALUE", "INHERITED"}, properties)
	return err
}
