package sinks

import (
	"io"
	"os"
	"path/filepath"
)

func ToFile(target os.DirEntry) (Sink, error) {
	var targetFile string
	var stream io.WriteCloser

	if target.IsDir() {
		targetFile = filepath.Join(target.Name(), "watchdog.log")
	} else {
		targetFile = target.Name()
	}

	_, err := os.Stat(targetFile)
	if os.IsNotExist(err) {
		stream, err = os.Create(targetFile)
		if err != nil {
			return nil, err
		}
	} else {
		stream, err = os.Open(targetFile)
		if err != nil {
			return nil, err
		}
	}

	return streamSinkImpl{stream}, nil
}
