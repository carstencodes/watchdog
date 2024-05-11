package sinks

import (
	"io"
)

type Sink interface {
	GetWriter() io.Writer
}

type streamSinkImpl struct {
	writer io.Writer
}

func (s streamSinkImpl) GetWriter() io.Writer {
	return s.writer
}

var allSinks []sinkFactory

type sinkFactory interface {
	Name() string
	SupportsSink(s string) bool
	CreateSink(s string) (Sink, error)
}
