package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func createLog() *log.Logger {
	var log = log.New(os.Stdout, "watchdog", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	return log
}

func toFile(logger *log.Logger, target os.DirEntry) error {
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
			return err
		}
	} else {
		stream, err = os.Open(targetFile)
		if err != nil {
			return err
		}
	}

	targetWriter := io.MultiWriter(os.Stdout, stream)
	logger.SetOutput(targetWriter)

	return nil
}
