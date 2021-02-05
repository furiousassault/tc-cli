package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	commandDescribe "github.com/furiousassault/tc-cli/pkg/commands/describe"
	commandList "github.com/furiousassault/tc-cli/pkg/commands/list"
	commandLog "github.com/furiousassault/tc-cli/pkg/commands/log"
	commandRun "github.com/furiousassault/tc-cli/pkg/commands/run"
	commandToken "github.com/furiousassault/tc-cli/pkg/commands/token"
	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/output"
	apiClient "github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

func main() {
	cmdRoot := createCmdRoot()
	configPath := cmdRoot.Flags().StringP("config-path", "c", "", "Path to configuration")
	err := cmdRoot.Flags().Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	config, err := configuration.ConfigFromYAML(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("configuration", config)


	api, err := apiClient.InitAPI(*config)
	if err != nil {
		log.Fatal("API init failed: ", err)
	}

	cmdRoot.AddCommand(
		commandLog.CreateCommandBuildLog(api.Logs, &output.StringPrinterStdout{}),
		commandToken.CreateCommandTreeToken(*config, api.Token),
		commandDescribe.CreateCommandTreeDescribe(api.Builds, output.Output()),
		commandList.CreateCommandTreeList(api.Projects, api.Builds, output.Output()),
		commandRun.CreateCommandBuildConfigurationRun(api.BuildQueue, output.Output()),
	)
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

func createCmdRoot() *cobra.Command {
	return &cobra.Command{
		SilenceUsage: true,
	}
}
