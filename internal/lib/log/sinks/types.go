package sinks

import (
	"io"
)

type SinkSetup interface {
	SetSinks(sinks ...[]Sink)
}

type Sink interface {
	GetWriter() io.Writer
}

type streamSinkImpl struct {
	writer io.Writer
}

func (s streamSinkImpl) GetWriter() io.Writer {
	return s.writer
}
