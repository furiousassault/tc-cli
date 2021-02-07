package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/output"

	commandDescribe "github.com/furiousassault/tc-cli/pkg/commands/describe"
	commandList "github.com/furiousassault/tc-cli/pkg/commands/list"
	commandLog "github.com/furiousassault/tc-cli/pkg/commands/log"
	commandRun "github.com/furiousassault/tc-cli/pkg/commands/run"
	commandToken "github.com/furiousassault/tc-cli/pkg/commands/token"

	apiClient "github.com/furiousassault/tc-cli/pkg/teamcity/api"
)

func main() {
	cmdRoot := createCommandRoot()
	configPath := cmdRoot.PersistentFlags().StringP(
		"config-path", "c", os.Getenv("TC_CLI_CONFIG_PATH"), "Path to configuration",
	)

	// This pre-parsing attempt returns error because doesn't see flags defined after its execution.
	// It's not clear how to parse args partially before main parsing/execution routine.
	// there should be another way to do it, without globals and such hacks. To fix later.
	_ = cmdRoot.PersistentFlags().Parse(os.Args[1:])

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
		commandLog.CreateCommandBuildLog(api.Logs, output.NewStringPrinterStdout()),
		commandToken.CreateCommandTreeToken(*config, api.Token),
		commandDescribe.CreateCommandTreeDescribe(api.Builds, output.NewBuildDescriptionWriter(os.Stdout)),
		commandList.CreateCommandTreeList(api.Projects, api.Builds, output.NewListWriter(os.Stdout)),
		commandRun.CreateCommandBuildConfigurationRun(api.BuildQueue, output.NewTriggerResultWriter(os.Stdout)),
	)

	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}

func createCommandRoot() *cobra.Command {
	return &cobra.Command{
		SilenceUsage: true,
	}
}