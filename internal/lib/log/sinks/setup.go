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
