package sinks

import (
	"io"
	"os"
)

func stdOut() (Sink, error) {
	return streamSinkImpl{
		io.Writer(os.Stdout),
	}, nil
}

type stdOutSinkFactory struct {
}

func (f stdOutSinkFactory) SupportsSink(s string) bool {
	return s == "stdout" || s == "1"
}

func (f stdOutSinkFactory) CreateSink(_ string) (Sink, error) {
	return stdOut()
}

func (f stdOutSinkFactory) String() string {
	return "stdout"
}

func (f stdOutSinkFactory) Name() string {
	return "stdout"
}

func init() {
	allSinks = append(allSinks, stdOutSinkFactory{})
}
