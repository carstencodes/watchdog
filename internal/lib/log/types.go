package log

import (
	"log"
)

type Logger interface {
	Printf(format string, v ...interface{})
	GetLog() *log.Logger
}

type Level string

var allLevels = []Level{
	Debug, Info, Warning, Error, Fatal,
}

var levelMap = map[string]Level{
	"debug":   Debug,
	"info":    Info,
	"warning": Warning,
	"error":   Error,
	"fatal":   Fatal,
}

const (
	Debug   Level = "[DEBUG]"
	Info    Level = "[INFO ]"
	Warning Level = "[WARN ]"
	Error   Level = "[ERROR]"
	Fatal   Level = "[FATAL]"
)

type Log interface {
	Debug() Logger
	Info() Logger
	Warning() Logger
	Error() Logger
	Fatal() Logger
}
