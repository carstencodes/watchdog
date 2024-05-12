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
	"io"
)

type Setup interface {
	WithSink(sink Sink, err error) Setup
	Build() (io.Writer, error)
}

type setupImpl struct {
	sinks []Sink
	err   error
}

func newSetup() Setup {
	return &setupImpl{[]Sink{}, nil}
}

func (s *setupImpl) WithSink(sink Sink, err error) Setup {
	if s.err != nil {
		return s
	}
	if err != nil {
		s.err = err
		return s
	}

	s.sinks = append(s.sinks, sink)
	return s
}

func (s *setupImpl) Build() (io.Writer, error) {
	if s.err != nil {
		return nil, s.err
	}

	writers := make([]io.Writer, len(s.sinks))
	for i, sink := range s.sinks {
		writers[i] = sink.GetWriter()
	}

	return io.MultiWriter(writers...), nil
}
