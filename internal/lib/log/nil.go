package log

import (
	"log"
)

type nilLogger struct {
}

func (n nilLogger) Printf(_ string, _ ...interface{}) {}

func newNilLogger() Logger {
	return &nilLogger{}
}

func (nilLogger) GetLog() *log.Logger {
	return nil
}
