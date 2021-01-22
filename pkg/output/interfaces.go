package output

type Outputter interface {
	List(headers []string, data [][]string)
}
