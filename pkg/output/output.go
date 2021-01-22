package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableOutputter struct {
	table *tablewriter.Table
}

func Output() *TableOutputter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	return &TableOutputter{table: table}
}

func (tw *TableOutputter) List(headers []string, data [][]string) {
	tw.table.SetHeader(headers)
	tw.table.AppendBulk(data)
	tw.table.Render()
}
