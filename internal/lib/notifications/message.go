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
