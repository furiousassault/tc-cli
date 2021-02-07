package output

import (
	"fmt"
	"io"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

var (
	tableHeadersTriggerResult           = []string{"BUILD_ID", "QUEUED_BY", "STATE"}
	tableHeadersTriggerResultProperties = []string{"TYPE", "KEY", "VALUE", "INHERITED"}
)

type TriggerResultWriter struct {
	t *TableWriter
}

func NewTriggerResultWriter(writer io.Writer) TriggerResultWriter {
	return TriggerResultWriter{
		t: NewTableWriter(writer),
	}
}

func (w TriggerResultWriter) WriteTriggerResult(result subapi.TriggerResult) {
	w.t.WriteTable(
		tableHeadersTriggerResult,
		[][]string{{fmt.Sprint(result.BuildID), result.TriggeredBy, result.BuildState}},
	)
	fmt.Println("\nProperties")

	properties := make([][]string, 0)

	for _, property := range result.Parameters.Items {
		properties = append(properties, []string{property.Type, property.Name, property.Value, fmt.Sprint(property.Inherited)})
	}

	w.t.WriteTable(tableHeadersTriggerResultProperties, properties)
}
