package output

import "fmt"

type StringPrinterStdout struct{}

func NewStringPrinterStdout() StringPrinterStdout {
	return StringPrinterStdout{}
}

func (sps StringPrinterStdout) PrintString(s string) {
	fmt.Println(s)
}
