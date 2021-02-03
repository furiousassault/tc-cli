package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

var cmdRoot = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := configuration.InitConfigFromYAML(); err != nil {
			return err
		}

		return api.InitAPI(configuration.GetConfig())
	},
	SilenceUsage: true,
}

func Execute() {
	cmdListBuild.Flags().IntVarP(
		&listBuildCount,
		"number",
		"n",
		0,
		"number of latest builds to return in list",
	)

	cmdRoot.PersistentFlags().StringVarP(
		&configuration.ConfigPath,
		"config", "c",
		os.Getenv("CONFIG_PATH"),
		"path to YAML config file",
	)

	cmdToken.AddCommand(cmdTokenRotate)
	cmdList.AddCommand(cmdListProject, cmdListBuildType, cmdListBuild)
	cmdDescribe.AddCommand(cmdDescribeProject, cmdDescribeBuildType, cmdDescribeBuild)

	cmdRoot.AddCommand(cmdList, cmdDescribe, cmdLog, cmdToken, cmdRun)

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
