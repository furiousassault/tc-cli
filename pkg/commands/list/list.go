package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/output"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

func CreateCommandTreeList(
	projectsAPI ProjectsAPI, buildsGetter BuildsGetter, outputter output.Outputter) *cobra.Command {
	cmdList := &cobra.Command{
		Use:   "list <subcommand>",
		Short: "list subcommand tree",
	}
	cmdListProject := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"project"},
		Short:   "list projects",
		Args:    cobra.NoArgs,
		RunE:    createHandlerListProjects(projectsAPI, outputter),
	}
	cmdListBuildType := &cobra.Command{
		Use:     "buildtypes [project_id]",
		Aliases: []string{"buildtype"},
		Short:   "list buildTypes of project specified by id",
		Args:    cobra.ExactArgs(1),
		RunE:    createHandlerListBuildTypes(projectsAPI, outputter),
	}
	cmdListBuild := &cobra.Command{
		Use:     "builds <buildtype_id>",
		Aliases: []string{"build"},
		Short:   "list builds of buildType specified by id",
		Args:    cobra.ExactArgs(1),
	}
	buildsCount := cmdListBuild.Flags().IntP(
		"number",
		"n",
		10,
		"number of latest builds to return in list",
	)
	cmdListBuild.RunE = createHandlerListBuilds(buildsGetter, outputter, *buildsCount)
	cmdList.AddCommand(cmdListProject, cmdListBuildType, cmdListBuild)

	return cmdList
}

type ProjectsAPI interface {
	GetList() (refs subapi.ProjectsReferences, err error)
	GetBuildTypesList(projectId string) (refs subapi.BuildTypeReferences, err error)
}

type BuildsGetter interface {
	GetBuildsByBuildConf(buildTypeID string, count int) (builds subapi.Builds, err error)
}

func createHandlerListProjects(
	projectsAPI ProjectsAPI, outputter output.Outputter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return listProjects(projectsAPI, outputter)
	}
}

func createHandlerListBuildTypes(
	projectsAPI ProjectsAPI, outputter output.Outputter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		return listBuildTypes(projectsAPI, outputter, projectID)
	}
}

func createHandlerListBuilds(
	buildsGetter BuildsGetter, outputter output.Outputter, count int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		return listBuilds(buildsGetter, outputter, projectID, count)
	}
}

func listProjects(projectsAPI ProjectsAPI, outputter output.Outputter) error {
	projects, err := projectsAPI.GetList()
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range projects.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name), entry.Description})
	}

	outputter.PrintTable([]string{"ID", "NAME", "DESCRIPTION"}, data)
	return nil
}

func listBuildTypes(projectsAPI ProjectsAPI, outputter output.Outputter, projectID string) error {
	buildTypes, err := projectsAPI.GetBuildTypesList(projectID)
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range buildTypes.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name)})
	}

	outputter.PrintTable([]string{"ID", "NAME"}, data)
	return nil
}

func listBuilds(buildsGetter BuildsGetter, outputter output.Outputter, buildTypeID string, buildsCount int) error {
	builds, err := buildsGetter.GetBuildsByBuildConf(buildTypeID, buildsCount)
	if err != nil {
		return err
	}

	data := make([][]string, 0)

	for _, entry := range builds.Items {
		data = append(data, []string{fmt.Sprint(entry.ID), entry.Number, entry.State, entry.Status})
	}

	outputter.PrintTable([]string{"id", "number", "state", "status"}, data)
	return nil
}
