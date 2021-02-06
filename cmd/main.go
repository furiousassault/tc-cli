package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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

	rootFlagset := pflag.NewFlagSet("some", pflag.ContinueOnError)
	configPath := rootFlagset.StringP(
		"config-path", "c", os.Getenv("TC_CLI_CONFIG_PATH"), "Path to configuration",
	)

	// it's not clear how to parse args partially before main parsing/execution routine
	// this pre-parsing attempt failing cause doesn't see flags defined after its execution
	// there should be another way to do it, without globals
	_ = rootFlagset.Parse(os.Args[1:])
	cmdRoot.PersistentFlags().AddFlagSet(rootFlagset)

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
