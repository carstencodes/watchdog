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
	"fmt"
)

type configurationImpl struct {
	providers                []Provider
	sectionDefinitions       map[string]interface{}
	providersLocked          bool
	sectionDefinitionsLocked bool
	parsed                   bool
}

func newConfiguration() *configurationImpl {
	return &configurationImpl{
		providers:                make([]Provider, 0),
		sectionDefinitions:       make(map[string]interface{}),
		providersLocked:          false,
		sectionDefinitionsLocked: false,
		parsed:                   false,
	}
}

func (c *configurationImpl) AddProvider(provider Provider) error {
	if c.providersLocked {
		return fmt.Errorf("cannot add provider %s to a locked collection", provider.Name())
	}

	c.providers = append(c.providers, provider)
	return nil
}

func (c *configurationImpl) BuildDefinitions() (DefinitionSetup, error) {
	c.providersLocked = true
	return c, nil
}

func (c *configurationImpl) AddConfigSectionDefinition(name string, v interface{}) error {
	if c.sectionDefinitionsLocked {
		return fmt.Errorf("cannot add section definition %s to a locked collection", name)
	}

	c.sectionDefinitions[name] = v
	return nil
}

func (c *configurationImpl) BuildConfiguration() (Configuration, error) {
	c.sectionDefinitionsLocked = true

	for _, provider := range c.providers {
		for key, value := range c.sectionDefinitions {
			err := provider.ReadConfigSectionDefinition(key, &value)
			if err != nil {
				return nil, fmt.Errorf("error reading section definition %s: %v", key, err)
			}
		}
	}

	return c, nil
}

func (c *configurationImpl) Parse() error {
	if c.parsed {
		return nil
	}

	for _, provider := range c.providers {
		err := provider.Parse()
		if err != nil {
			return fmt.Errorf("error parsing provider %s: %v", provider.Name(), err)
		}
	}

	c.parsed = true

	return nil
}

func (c *configurationImpl) Parsed() bool {
	return c.parsed
}
