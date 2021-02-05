package describe

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

func CreateCommandTreeDescribe(buildGetter BuildGetter, outputter output.Outputter) *cobra.Command {
	cmdDescribe := &cobra.Command{
		Use:   "describe <subcommand>",
		Short: "describe subcommand tree",
	}
	cmdDescribeBuild := &cobra.Command{
		Use:   "build <build_typeId> <build_number>",
		Short: "show main attributes of entity specified by id",
		Args:  cobra.ExactArgs(2),
		RunE:  createHandlerDescribeBuild(buildGetter, outputter),
	}
	cmdDescribe.AddCommand(cmdDescribeBuild)

	return cmdDescribe
}

func createHandlerDescribeBuild(buildGetter BuildGetter, outputter output.Outputter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		buildTypeID := args[0]
		buildNumber := args[1]

		return describeBuild(buildGetter, outputter, buildTypeID, buildNumber)
	}
}

type BuildGetter interface {
	GetBuild(buildTypeID string, number string) (build subapi.Build, err error)
}

func describeBuild(buildGetter BuildGetter, outputter output.Outputter, buildTypeID string, number string) error {
	build, err := buildGetter.GetBuild(buildTypeID, number)
	if err != nil {
		return err
	}

	outputter.PrintTable(
		[]string{"ID", "NUMBER", "STATE", "STATUS", "QUEUED", "STARTED", "FINISHED"},
		[][]string{
			{
				fmt.Sprint(build.ID),
				build.Number,
				build.State,
				build.Status,
				build.QueuedDate,
				build.StartDate,
				build.FinishDate,
			},
		},
	)
	fmt.Println()
	fmt.Println("Properties")

	properties := make([][]string, 0)

	for _, property := range build.Properties.Items {
		properties = append(properties, []string{property.Name, property.Value, fmt.Sprint(property.Inherited)})
	}

	outputter.PrintTable([]string{"KEY", "VALUE", "INHERITED"}, properties)
	return nil
}
