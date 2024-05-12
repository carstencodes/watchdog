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

package sinks

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type sinkFactoryParam struct {
	parameter string
	factory   sinkFactory
}

func (s sinkFactoryParam) Invoke() (Sink, error) {
	return s.factory.CreateSink(s.parameter)
}

func (s sinkFactoryParam) String() string {
	return s.factory.Name()
}

type sinksCollectionVar []sinkFactoryParam

func (s *sinksCollectionVar) String() string {
	builder := strings.Builder{}
	hasValues := false
	for _, param := range *s {
		if hasValues {
			builder.WriteString(", ")
		}
		builder.WriteString(param.String())
		hasValues = true
	}
	return builder.String()
}

func (s *sinksCollectionVar) Set(value string) error {
	for _, factory := range allSinks {
		if factory.SupportsSink(value) {
			*s = append(*s, sinkFactoryParam{value, factory})
			return nil
		}
	}

	return fmt.Errorf("invalid sink to obtain: %s", value)
}

var sinksVar = sinksCollectionVar{sinkFactoryParam{"", stdOutSinkFactory{}}}

func init() {
	flag.Var(&sinksVar, "log-to", "Select the logging parts to write to. Can be applied multiple times")
}

func CreateSink() (io.Writer, error) {
	setup := newSetup()

	for _, factory := range sinksVar {
		s, err := factory.Invoke()
		setup.WithSink(s, err)
	}

	return setup.Build()
}
