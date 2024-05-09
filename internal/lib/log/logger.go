package log

import (
	"io"
	"log"
)

type loggerImpl struct {
	logger *log.Logger
}

func newLogger(level Level, writer io.Writer) Logger {
	lg := log.New(writer, "watchdog "+string(level)+": ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile|log.Lmsgprefix)

	if level == Fatal {
		return newFatalLogger(lg)
	}

	return &loggerImpl{lg}
}

func (l loggerImpl) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l loggerImpl) GetLog() *log.Logger {
	return l.logger
}
