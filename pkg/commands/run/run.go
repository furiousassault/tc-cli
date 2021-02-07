package run

import (
	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

type BuildRunner interface {
	RunBuildByBuildConfID(buildconfID string) (result subapi.TriggerResult, err error)
}

type triggerResultWriter interface {
	WriteTriggerResult(result subapi.TriggerResult)
}

func CreateCommandBuildConfigurationRun(buildRunner BuildRunner, writer triggerResultWriter) *cobra.Command {
	return &cobra.Command{
		Use:   "run <build_configuration_id>",
		Short: "run build configuration",
		Args:  cobra.ExactArgs(1),
		RunE:  createHandlerBuildConfigurationRun(buildRunner, writer),
	}
}

func createHandlerBuildConfigurationRun(
	buildRunner BuildRunner, writer triggerResultWriter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		buildConfID := args[0]

		return buildConfigurationRun(buildRunner, writer, buildConfID)
	}

}

func buildConfigurationRun(buildRunner BuildRunner, writer triggerResultWriter, buildConfID string) error {
	result, err := buildRunner.RunBuildByBuildConfID(buildConfID)
	if err != nil {
		return err
	}

	writer.WriteTriggerResult(result)
	return err
}
