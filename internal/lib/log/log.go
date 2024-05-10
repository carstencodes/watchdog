package log

import (
	"github.com/carstencodes/watchdog/internal/lib/common"
)

func CreateLog(minLevel Level, info common.ApplicationDetails, setup Setup) (Log, error) {
	levels := getLogLevel(minLevel)
	writer, err := setup.Build()
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
