package output

import (
	"fmt"
	"io"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

var (
	tableHeadersBuildDescription         = []string{"ID", "NUMBER", "STATE", "STATUS", "QUEUED", "STARTED", "FINISHED"}
	tableHeadersBuildPropertiesInput     = []string{"KEY", "VALUE", "INHERITED"}
	tableHeadersBuildPropertiesResulting = []string{"KEY", "VALUE"}
)

type BuildDescriptionWriter struct {
	t *TableWriter
}

func NewBuildDescriptionWriter(writer io.Writer) *BuildDescriptionWriter {
	return &BuildDescriptionWriter{
		t: NewTableWriter(writer),
	}
}

func (w *BuildDescriptionWriter) WriteBuildDescription(build subapi.BuildJson) {
	w.t.WriteTable(
		tableHeadersBuildDescription,
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
	fmt.Println("\nProperties")
	properties := make([][]string, 0)
	resultingProperties := make([][]string, 0)

	for _, property := range build.Properties.Items {
		properties = append(
			properties,
			[]string{property.Name, property.Value, fmt.Sprint(property.Inherited)},
		)
	}

	w.t.WriteTable(tableHeadersBuildPropertiesInput, properties)

	for _, property := range build.ResultingProperties.Items {
		resultingProperties = append(
			resultingProperties,
			[]string{property.Name, property.Value},
		)
	}

	if len(resultingProperties) > 0 {
		fmt.Println("\nResulting properties")
		w.t.WriteTable(tableHeadersBuildPropertiesResulting, resultingProperties)
	}
}
