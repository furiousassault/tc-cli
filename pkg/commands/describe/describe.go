package describe

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

type buildGetter interface {
	GetBuild(buildTypeID string, number string) (build subapi.BuildJson, err error)
	GetBuildResults(buildID string) (resultingProperties subapi.Properties, err error)
}

type buildDescriptionWriter interface {
	WriteBuildDescription(build subapi.BuildJson)
}

func CreateCommandTreeDescribe(buildGetter buildGetter, writer buildDescriptionWriter) *cobra.Command {
	cmdDescribe := &cobra.Command{
		Use:   "describe <subcommand>",
		Short: "describe subcommand tree",
	}
	cmdDescribeBuild := &cobra.Command{
		Use:   "build <buildtype_ID> <build_number>",
		Short: "Show main attributes of build specified by id",
		Args:  cobra.ExactArgs(2),
	}
	resultingPropertiesPointer := cmdDescribeBuild.Flags().BoolP(
		"short",
		"s",
		false,
		"don't output resulting properties of build",
	)
	cmdDescribeBuild.RunE = createHandlerDescribeBuild(buildGetter, writer, resultingPropertiesPointer)
	cmdDescribe.AddCommand(cmdDescribeBuild)

	return cmdDescribe
}

func createHandlerDescribeBuild(
	buildGetter buildGetter,
	writer buildDescriptionWriter,
	shortFlag *bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		buildTypeID := args[0]
		buildNumber := args[1]

		return describeBuild(buildGetter, writer, buildTypeID, buildNumber, *shortFlag)
	}
}

func describeBuild(
	buildGetter buildGetter, writer buildDescriptionWriter,
	buildTypeID string, number string, shortFlag bool) error {
	build, err := buildGetter.GetBuild(buildTypeID, number)
	if err != nil {
		return err
	}

	if !shortFlag {
		build.ResultingProperties, err = buildGetter.GetBuildResults(fmt.Sprint(build.ID))
		if err != nil {
			return err
		}
	}

	writer.WriteBuildDescription(build)

	return nil
}
