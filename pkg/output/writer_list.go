package output

import (
	"fmt"
	"io"

	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

var (
	tableHeadersListProjects   = []string{"ID", "NAME", "DESCRIPTION"}
	tableHeadersListBuildTypes = []string{"ID", "NAME"}
	tableHeadersListBuilds     = []string{"ID", "NUMBER", "STATE", "STATUS"}
)

type ListWriter struct {
	t *TableWriter
}

func NewListWriter(writer io.Writer) ListWriter {
	return ListWriter{
		t: NewTableWriter(writer),
	}
}

func (w ListWriter) WriteListProjects(projects subapi.ProjectsReferences) {
	data := make([][]string, 0)

	for _, entry := range projects.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name), entry.Description})
	}

	w.t.WriteTable(tableHeadersListProjects, data)
}

func (w ListWriter) WriteListBuildTypes(buildTypes subapi.BuildTypeReferences) {
	data := make([][]string, 0)

	for _, entry := range buildTypes.Items {
		data = append(data, []string{entry.ID, fmt.Sprintf("\"%s\"", entry.Name)})
	}

	w.t.WriteTable(tableHeadersListBuildTypes, data)
}

func (w ListWriter) WriteListBuilds(builds subapi.Builds) {
	data := make([][]string, 0)

	for _, entry := range builds.Items {
		data = append(data, []string{fmt.Sprint(entry.ID), entry.Number, entry.State, entry.Status})
	}

	w.t.WriteTable(tableHeadersListBuilds, data)
}
