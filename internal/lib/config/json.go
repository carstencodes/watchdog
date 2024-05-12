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

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type jsonFileConfigProvider struct {
	searchPaths    []string
	configSections map[string]interface{}
	fileMustExist  bool
}

func NewJsonFileConfigProvider(fileMustExist bool) FileProvider {
	return &jsonFileConfigProvider{
		searchPaths:    make([]string, 0),
		configSections: make(map[string]interface{}),
		fileMustExist:  fileMustExist,
	}
}

func (j *jsonFileConfigProvider) AddSearchPath(path string) error {
	if len(path) == 0 {
		return nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf(`path "%s" is not a directory`, path)
	}

	j.searchPaths = append(j.searchPaths, path)

	return nil
}

func (j *jsonFileConfigProvider) Name() string {
	return "json"
}

func (j *jsonFileConfigProvider) ReadConfigSectionDefinition(name string, v interface{}) error {
	_, ok := j.configSections[name]
	if !ok {
		j.configSections[name] = v
		return nil
	}

	return fmt.Errorf("config section '%s' already exists", name)
}

func (j *jsonFileConfigProvider) Parse() error {
	if strings.HasSuffix(configFilePath, ".json") {
		filePath := ""
		if _, err := os.Stat(configFilePath); err == nil {
			filePath = configFilePath
		} else if !os.IsNotExist(err) {
			return err
		} else {
			for _, searchPath := range j.searchPaths {
				configFile := filepath.Join(searchPath, configFilePath)
				_, err := os.Stat(configFile)
				if err != nil {
					if !os.IsNotExist(err) {
						return err
					}

					continue
				}

				filePath = configFile
				break
			}
		}

		if len(filePath) == 0 {
			if j.fileMustExist {
				return fmt.Errorf("config file %s not found in search path: %s", configFilePath, j.searchPaths)
			}

			return nil
		}

		var err error
		var innerError error
		err = j.loadJson(filePath, func(e error) { innerError = e })
		if err != nil {
			return err
		}
		if innerError != nil {
			return innerError
		}
	}
	return nil
}

func (j *jsonFileConfigProvider) loadJson(filePath string, handlerFunc func(error)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			handlerFunc(err)
		}
	}()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&j.configSections)
}
