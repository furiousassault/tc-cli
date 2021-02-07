package output

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	w io.Writer
}

func NewTableWriter(writer io.Writer) *TableWriter {
	return &TableWriter{w: writer}
}

func (tw *TableWriter) WriteTable(headers []string, data [][]string) {
	// olekukonko/tablewriter table instance has state, so in our case it's easier to instantiate another instance
	// each time instead of resetting and cleaning the old one.
	table := tablewriter.NewWriter(tw.w)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.SetHeader(headers)
	table.AppendBulk(data)
	table.Render()
}
