package sinks

import (
	"io"
	"os"
)

func StdErr() (Sink, error) {
	return streamSinkImpl{
		io.Writer(os.Stderr),
	}, nil
}
