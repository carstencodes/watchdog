package sinks

import (
	"io"
	"os"
)

func StdOut() (Sink, error) {
	return streamSinkImpl{
		io.Writer(os.Stdout),
	}, nil
}
