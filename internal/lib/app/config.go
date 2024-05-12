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

	"github.com/carstencodes/watchdog/internal/lib/common"
	"github.com/carstencodes/watchdog/internal/lib/config"
)

func setupConfiguration() error {
	var err error = nil
	config.Sections, err = config.Providers.BuildDefinitions()
	if err != nil {
		return err
	}

	config.Values, err = config.Sections.BuildConfiguration()
	if err != nil {
		return err
	}

	if !config.Values.Parsed() {
		return config.Values.Parse()
	}

	return nil
}

func init() {
	err := initConfig()
	if err != nil {
		panic(err)
	}
}

func initConfig() error {
	var err error = nil
	var namespace = common.ApplicationInfo().Name()

	json := config.NewJsonFileConfigProvider(false)
	err = initOs(json)
	if err != nil {
		return err
	}
	env := config.NewEnvProvider(namespace)
	flags := config.NewFlagConfigProvider()

	err = config.Providers.AddProvider(json)
	if err != nil {
		return err
	}
	err = config.Providers.AddProvider(env)
	if err != nil {
		return err
	}
	err = config.Providers.AddProvider(flags)

	return err
}

func getEnvDir(variableName string) string {
	value := os.Getenv(variableName)
	if value == "" {
		return ""
	}

	return filepath.Join(value, common.ApplicationInfo().Name())
}
