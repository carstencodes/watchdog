// Copyright (c) 2022-2024 Carsten Igel
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
