package log

import (
	"log"
)

type fatalLogger struct {
	log *log.Logger
}

func newFatalLogger(logger *log.Logger) Logger {
	return fatalLogger{log: logger}
}

func (f fatalLogger) Printf(format string, args ...interface{}) {
	f.log.Fatalf(format, args...)
}

func (f fatalLogger) GetLog() *log.Logger {
	return f.log
}
