// Copyright (c) 2022-2024 Carsten Igel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
