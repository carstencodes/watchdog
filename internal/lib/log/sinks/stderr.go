package sinks

import (
	"io"
	"os"
)

func stdErr() (Sink, error) {
	return streamSinkImpl{
		io.Writer(os.Stderr),
	}, nil
}

type stdErrSinkFactory struct {
}

func (f stdErrSinkFactory) SupportsSink(s string) bool {
	return s == "stderr" || s == "2"
}

func (f stdErrSinkFactory) CreateSink(_ string) (Sink, error) {
	return stdErr()
}

func (f stdErrSinkFactory) String() string {
	return "stderr"
}

func (f stdErrSinkFactory) Name() string {
	return "stderr"
}

func init() {
	allSinks = append(allSinks, stdErrSinkFactory{})
}
