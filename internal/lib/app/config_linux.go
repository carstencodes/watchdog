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

package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/carstencodes/watchdog/internal/lib/common"
	"github.com/carstencodes/watchdog/internal/lib/config"
)

func initOs(fileProvider config.FileProvider) error {
	var err error = nil
	err = fileProvider.AddSearchPath("/etc")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/etc/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/etc/opt/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath("/opt/etc/" + common.ApplicationInfo().Name())
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/." + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/.config/" + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}
	err = fileProvider.AddSearchPath(os.Getenv("HOME") + "/.local/etc/" + common.ApplicationInfo().Name() + "/")
	if err != nil {
		return err
	}

	xdgDataDirs := os.Getenv("XDG_DATA_DIRS")
	xdgDataDirValues := strings.Split(xdgDataDirs, ":")
	for _, xdgDataDir := range xdgDataDirValues {
		err = fileProvider.AddSearchPath(filepath.Join(xdgDataDir, common.ApplicationInfo().Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
