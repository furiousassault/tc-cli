package run

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

func CreateCommandBuildConfigurationRun(buildRunner BuildRunner, outputter output.Outputter) *cobra.Command {
	return &cobra.Command{
		Use:   "run <build_configuration_id>",
		Short: "run build configuration",
		Args:  cobra.ExactArgs(1),
		RunE:  createHandlerBuildConfigurationRun(buildRunner, outputter),
	}
}

func createHandlerBuildConfigurationRun(
	buildRunner BuildRunner, outputter output.Outputter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		buildConfID := args[0]

		return buildConfigurationRun(buildRunner, outputter, buildConfID)
	}

}

type BuildRunner interface {
	RunBuildByBuildConfID(buildconfID string) (result subapi.TriggerResult, err error)
}

func buildConfigurationRun(buildRunner BuildRunner, outputter output.Outputter, buildConfID string) error {
	result, err := buildRunner.RunBuildByBuildConfID(buildConfID)
	if err != nil {
		return err
	}

	outputter.PrintTable(
		[]string{"BUILD_ID", "QUEUED_BY", "STATE"},
		[][]string{{fmt.Sprint(result.BuildID), result.TriggeredBy, result.BuildState}},
	)
	fmt.Println()
	fmt.Println("Properties")

	properties := make([][]string, 0)

	for _, property := range result.Parameters.Items {
		properties = append(properties, []string{property.Type, property.Name, property.Value, fmt.Sprint(property.Inherited)})
	}

	outputter.PrintTable([]string{"TYPE", "KEY", "VALUE", "INHERITED"}, properties)
	return err
}
