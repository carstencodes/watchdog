package log

import (
	"log"
)

type nilLogger struct {
}

func (n nilLogger) Printf(format string, args ...interface{}) {}

func NewNilLogger() Logger {
	return &nilLogger{}
}

func (nilLogger) GetLog() *log.Logger {
	return nil
}
