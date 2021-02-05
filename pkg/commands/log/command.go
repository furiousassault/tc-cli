package log

import (
	"github.com/spf13/cobra"
)

type BuildLogGetter interface {
	GetBuildLog(buildID string) (log []byte, err error)
}

type StringPrinter interface {
	PrintString(s string)
}

func CreateCommandBuildLog(getter BuildLogGetter, printer StringPrinter) *cobra.Command {
	return &cobra.Command{
		Use:   "log <build id>",
		Short: "Show logs for of the build specified by id",
		Args:  validateLogArgs,
		RunE:  createHandler(getter, printer),
	}
}

func createHandler(getter BuildLogGetter, printer StringPrinter) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return buildLog(getter, printer, args[0])
	}
}

func buildLog(getter BuildLogGetter, printer StringPrinter, buildID string) error {
	out, err := getter.GetBuildLog(buildID)
	if err != nil {
		return err
	}

	printer.PrintString(string(out))
	return nil
}
