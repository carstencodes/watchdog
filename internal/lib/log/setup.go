package log

import (
	"io"

	"github.com/carstencodes/watchdog/internal/lib/log/sinks"
)

type Setup interface {
	WithSink(sink sinks.Sink, err error) Setup
	Build() (io.Writer, error)
}

type setupImpl struct {
	sinks []sinks.Sink
	err   error
}

func NewSetup() Setup {
	return &setupImpl{[]sinks.Sink{}, nil}
}

func (s *setupImpl) WithSink(sink sinks.Sink, err error) Setup {
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
