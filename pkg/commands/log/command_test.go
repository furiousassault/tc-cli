package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

type buildLogGetterMock struct {
	log []byte
	err error
}

func newBuildLogGetterMock(log []byte, err error) *buildLogGetterMock {
	return &buildLogGetterMock{log: log, err: err}
}

func (m *buildLogGetterMock) GetBuildLog(_ string) (log []byte, err error) {
	return m.log, m.err
}

type stringPrinterMock struct {
	buffer []string
}

func newStringPrinterMock() *stringPrinterMock {
	return &stringPrinterMock{buffer: make([]string, 0)}
}

func (s *stringPrinterMock) PrintString(str string) {
	s.buffer = append(s.buffer, str)
}

func (s *stringPrinterMock) last() string {
	if len(s.buffer) < 1 {
		return ""
	}

	return s.buffer[len(s.buffer)-1]
}

func TestBuildLog_Negative(t *testing.T) {
	getter := newBuildLogGetterMock([]byte("test"), errTest)
	printer := newStringPrinterMock()
	err := buildLog(getter, printer, "")
	assert.Error(t, err)
	assert.Equal(t, "", printer.last())
}

func TestBuildLog_Positive(t *testing.T) {
	getter := newBuildLogGetterMock([]byte("test"), nil)
	printer := newStringPrinterMock()
	err := buildLog(getter, printer, "")
	assert.NoError(t, err)
	assert.Equal(t, "test", printer.last())
}
