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
