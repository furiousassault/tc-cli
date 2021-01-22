package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

var (
	cmdList = &cobra.Command{
		Use:   "list <subcommand>",
		Short: "list subcommand tree",
	}

	cmdListProject = &cobra.Command{
		Use:     "projects",
		Aliases: []string{"project"},
		Short:   "list projects",
		Args:    cobra.NoArgs,
		RunE:    listProjects,
	}

	cmdListBuildType = &cobra.Command{
		Use:     "buildtypes [project_id]",
		Aliases: []string{"buildtype"},
		Short:   "list buildTypes of project specified by id",
		Args:    cobra.ExactArgs(1),
		RunE:    listBuildTypes,
	}

	cmdListBuild = &cobra.Command{
		Use:     "builds <buildtype_id>",
		Aliases: []string{"build"},
		Short:   "list builds of buildType specified by id",
		Args:    cobra.ExactArgs(1),
		RunE:    listBuilds,
	}

	listBuildCount int
)

func listProjects(_ *cobra.Command, args []string) error {
	projects, err := api.API().Projects.GetList()
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range projects.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name), entry.Description})
	}

	output.Output().List([]string{"ID", "NAME", "DESCRIPTION"}, data)
	return nil
}

func listBuildTypes(_ *cobra.Command, args []string) error {
	buildTypes, err := api.API().Projects.GetBuildTypesList(args[0])
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range buildTypes.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name)})
	}

	output.Output().List([]string{"ID", "NAME"}, data)
	return nil
}

func listBuilds(_ *cobra.Command, args []string) error {
	builds, err := api.API().Builds.GetBuildsByBuildConf(args[0], listBuildCount)
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range builds.Items {
		data = append(data, []string{fmt.Sprint(entry.ID), entry.Number, entry.State, entry.Status})
	}

	output.Output().List([]string{"id", "number", "state", "status"}, data)
	return nil
}
