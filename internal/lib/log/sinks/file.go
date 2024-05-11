package sinks

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func toFile(target os.DirEntry) (Sink, error) {
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

type fileSinkFactory struct {
}

func (f fileSinkFactory) SupportsSink(s string) bool {
	return strings.HasPrefix(s, "f:")
}

func (f fileSinkFactory) CreateSink(s string) (Sink, error) {
	path := s[2 : len(s)-1]
	if !fs.ValidPath(path) {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	entry := filePathEntry{stat}

	return toFile(entry)
}

func (f fileSinkFactory) String() string {
	return "f:<PathToLogFileOrDir>"
}

func (f fileSinkFactory) Name() string {
	return "file"
}

func init() {
	allSinks = append(allSinks, fileSinkFactory{})
}

type filePathEntry struct {
	filePath fs.FileInfo
}

func (f filePathEntry) Name() string {
	return f.filePath.Name()
}

func (f filePathEntry) IsDir() bool {
	return f.filePath.IsDir()
}

func (f filePathEntry) Type() fs.FileMode {
	return f.filePath.Mode()
}

func (f filePathEntry) Info() (fs.FileInfo, error) {
	return f.filePath, nil
}
