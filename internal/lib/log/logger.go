package log

import (
	"fmt"
	"io"
	"log"
)

type loggerImpl struct {
	logger *log.Logger
}

func newLogger(level Level, applicationName string, writer io.Writer) Logger {
	lg := log.New(writer, applicationName+" "+string(level)+": ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile|log.Lmsgprefix)

	if level == Fatal {
		return newFatalLogger(lg)
	}

	return &loggerImpl{lg}
}

func (l loggerImpl) Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	err := l.logger.Output(2, message)
	if err != nil {
		l.logger.Printf(format, args...) // log anyway - ignore error
	}
}

func (l loggerImpl) GetLog() *log.Logger {
	return l.logger
}
