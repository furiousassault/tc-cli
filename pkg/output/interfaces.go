package output

type Outputter interface {
	PrintTable(headers []string, data [][]string)
}
