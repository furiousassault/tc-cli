package output

import "fmt"

type StringPrinterStdout struct{}

func (sps *StringPrinterStdout) PrintString(s string) {
	fmt.Println(s)
}
