package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

var (
	cmdDescribe = &cobra.Command{
		Use:   "describe <subcommand>",
		Short: "describe subcommand tree",
		RunE:  nil,
	}
	cmdDescribeProject = &cobra.Command{
		Use:   "project <project_id>",
		Short: "show main attributes of entity specified by id",
		RunE:  nil,
	}
	cmdDescribeBuildType = &cobra.Command{
		Use:   "buildtype <buildtype_id>",
		Short: "show main attributes of entity specified by id",
		RunE:  nil,
	}
	cmdDescribeBuild = &cobra.Command{
		Use:   "build <build_typeId> <build_number>",
		Short: "show main attributes of entity specified by id",
		Args:  cobra.ExactArgs(2),
		RunE:  describeBuild,
	}
)

func describeBuild(_ *cobra.Command, args []string) error {
	build, err := api.API().Builds.GetBuild(args[0], args[1])
	if err != nil {
		return err
	}

	output.Output().List(
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

	output.Output().List([]string{"KEY", "VALUE", "INHERITED"}, properties)
	return nil
}
