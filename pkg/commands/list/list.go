package list

import (
	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

type projectsAPI interface {
	GetList() (refs subapi.ProjectsReferences, err error)
	GetBuildTypesList(projectId string) (refs subapi.BuildTypeReferences, err error)
}

type buildsGetter interface {
	GetBuildsByBuildConf(buildTypeID string, count int) (builds subapi.Builds, err error)
}

type listWriter interface {
	WriteListProjects(projects subapi.ProjectsReferences)
	WriteListBuildTypes(buildTypes subapi.BuildTypeReferences)
	WriteListBuilds(builds subapi.Builds)
}

func CreateCommandTreeList(
	projectsAPI projectsAPI, buildsGetter buildsGetter, writer listWriter) *cobra.Command {
	cmdList := &cobra.Command{
		Use:   "list <subcommand>",
		Short: "list subcommand tree",
	}
	cmdListProject := &cobra.Command{
		Use:     "projects",
		Aliases: []string{"project"},
		Short:   "list projects",
		Args:    cobra.NoArgs,
		RunE:    createHandlerListProjects(projectsAPI, writer),
	}
	cmdListBuildType := &cobra.Command{
		Use:     "buildtypes [project_id]",
		Aliases: []string{"buildtype"},
		Short:   "list buildTypes of project specified by id",
		Args:    cobra.ExactArgs(1),
		RunE:    createHandlerListBuildTypes(projectsAPI, writer),
	}
	cmdListBuild := &cobra.Command{
		Use:     "builds <buildtype_id>",
		Aliases: []string{"build"},
		Short:   "list builds of buildType specified by id",
		Args:    cobra.ExactArgs(1),
	}

	buildsCountPointer := cmdListBuild.Flags().IntP(
		"number",
		"n",
		10,
		"number of latest builds to return in list",
	)
	// pass flag as a pointer to make it filled with actual value on final function execution;
	// todo find more obvious solution to create command flags
	cmdListBuild.RunE = createHandlerListBuilds(buildsGetter, writer, buildsCountPointer)
	cmdList.AddCommand(cmdListProject, cmdListBuildType, cmdListBuild)

	return cmdList
}

func createHandlerListProjects(
	projectsAPI projectsAPI, writer listWriter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return listProjects(projectsAPI, writer)
	}
}

func createHandlerListBuildTypes(
	projectsAPI projectsAPI, writer listWriter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		return listBuildTypes(projectsAPI, writer, projectID)
	}
}

func createHandlerListBuilds(
	buildsGetter buildsGetter, writer listWriter, count *int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		return listBuilds(buildsGetter, writer, projectID, *count)
	}
}

func listProjects(projectsAPI projectsAPI, writer listWriter) error {
	projects, err := projectsAPI.GetList()
	if err != nil {
		return err
	}

	writer.WriteListProjects(projects)
	return nil
}

func listBuildTypes(projectsAPI projectsAPI, writer listWriter, projectID string) error {
	buildTypes, err := projectsAPI.GetBuildTypesList(projectID)
	if err != nil {
		return err
	}

	writer.WriteListBuildTypes(buildTypes)
	return nil
}

func listBuilds(buildsGetter buildsGetter, writer listWriter, buildTypeID string, buildsCount int) error {
	builds, err := buildsGetter.GetBuildsByBuildConf(buildTypeID, buildsCount)
	if err != nil {
		return err
	}

	writer.WriteListBuilds(builds)
	return nil
}
