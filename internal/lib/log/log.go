package log

import (
	"flag"
	"fmt"
	"strings"

	"github.com/carstencodes/watchdog/internal/lib/common"
	"github.com/carstencodes/watchdog/internal/lib/log/sinks"
)

var level = logLevelVar{Info}

type logLevelVar struct {
	levelValue Level
}

func (l logLevelVar) String() string {
	for key, value := range levelMap {
		if value == l.levelValue {
			return key
		}
	}

	return string(l.levelValue)
}

func (l logLevelVar) Set(value string) error {
	value = strings.ToLower(value)
	if newLevel, ok := levelMap[value]; ok {
		l.levelValue = newLevel
		return nil
	} else {
		return fmt.Errorf("unknown log level '%s'", value)
	}
}

func init() {
	flag.Var(&level, "log-level", "Select the log level to use. Must be one of: debug, info, warning, error, fatal")
}

func CreateLog(info common.ApplicationDetails) (Log, error) {
	minLevel := level.levelValue
	levels := getLogLevel(minLevel)
	writer, err := sinks.CreateSink()
	if err != nil {
		return nil, err
	}

	return newLogShell(levels, info, writer), nil
}

func getLogLevel(minLevel Level) []Level {
	var levels []Level
	switch minLevel {
	case Debug:
		levels = append(levels, Debug)
		fallthrough
	case Info:
		levels = append(levels, Info)
		fallthrough
	case Warning:
		levels = append(levels, Warning)
		fallthrough
	case Error:
		levels = append(levels, Error)
		fallthrough
	case Fatal:
		levels = append(levels, Fatal)
	}
	return levels
}
