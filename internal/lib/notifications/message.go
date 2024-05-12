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

package notifications

import (
	"errors"
	"fmt"
	"strconv"
)

type argsMap struct {
	items map[string]string
}

type Message interface {
	getTitle() string
	getMessage() string
}

type messageImpl struct {
	title   string
	message string
}

func NewArgsMap(items map[string]string) Args {
	return argsMap{items: items}
}

func (m messageImpl) getTitle() string {
	return m.title
}

func (m messageImpl) getMessage() string {
	return m.message
}

func NewMessage(title string, message string) Message {
	return messageImpl{title: title, message: message}
}

func (am argsMap) GetBool(key string) (bool, error) {
	val, found := am.items[key]
	if !found {
		return false, errors.New(fmt.Sprintf("Failed to find key %s", key))
	}

	value, err := strconv.ParseBool(val)
	return value, err
}

func (am argsMap) GetFloat(key string) (float64, error) {
	val, found := am.items[key]
	if !found {
		return 0.0, errors.New(fmt.Sprintf("Failed to find key %s", key))
	}

	value, err := strconv.ParseFloat(val, 64)
	return value, err
}

func (am argsMap) GetInt(key string) (int64, error) {
	val, found := am.items[key]
	if !found {
		return 0, errors.New(fmt.Sprintf("Failed to find key %s", key))
	}

	value, err := strconv.ParseInt(val, 10, 64)
	return value, err
}

func (am argsMap) GetString(key string) (string, error) {
	val, found := am.items[key]
	if !found {
		return "", errors.New(fmt.Sprintf("Failed to find key %s", key))
	}

	return val, nil
}
