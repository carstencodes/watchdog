package log

import (
	"io"

	"github.com/carstencodes/watchdog/internal/lib/common"
)

type logShell struct {
	loggers map[Level]Logger
}

func newLogShell(levels []Level, info common.ApplicationDetails, writer io.Writer) logShell {
	loggers := make(map[Level]Logger)
	for _, level := range levels {
		loggers[level] = newLogger(level, info.Name(), writer)
	}

	for _, level := range allLevels {
		if _, found := loggers[level]; !found {
			loggers[level] = newNilLogger()
		}
	}

	return logShell{loggers: loggers}
}

func (l logShell) Debug() Logger {
	if l, ok := l.loggers[Debug]; ok {
		return l
	}

	return newNilLogger()
}

func (l logShell) Info() Logger {
	if l, ok := l.loggers[Info]; ok {
		return l
	}

	return newNilLogger()
}

func (l logShell) Warning() Logger {
	if l, ok := l.loggers[Warning]; ok {
		return l
	}

	return newNilLogger()
}

func (l logShell) Error() Logger {
	if l, ok := l.loggers[Error]; ok {
		return l
	}

	return newNilLogger()
}

func (l logShell) Fatal() Logger {
	if l, ok := l.loggers[Fatal]; ok {
		return l
	}

	return newNilLogger()
}
