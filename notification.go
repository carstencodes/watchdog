package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type Args interface {
	GetBool(key string) (bool, error)
	GetFloat(key string) (float64, error)
	GetInt(key string) (int64, error)
	GetString(key string) (string, error)
}

type Notifier interface {
	Connect() error
	Disconnect() error
	Send(message Message, messageArgs Args) error
}

type Message struct {
	title   string
	message string
}

type NotifierCreatorFunc func() (Notifier, error)

var notification_client string
var notification_clients map[string]NotifierCreatorFunc = make(map[string]NotifierCreatorFunc)

func init() {
	flag.StringVar(&notification_client, "notify", "", "The notification client")

	notification_clients[""] = newNullNotifier
}

func getNotificationClient() (Notifier, error) {
	creator, present := notification_clients[notification_client]

	if !present {
		return nil, errors.New(fmt.Sprintf("Unknown notification client: %s", notification_client))
	}

	return creator()
}

type nullNotifier struct{}

func newNullNotifier() (Notifier, error) {
	return nullNotifier{}, nil
}

func (nn nullNotifier) Connect() error {
	return nil
}

func (nn nullNotifier) Disconnect() error {
	return nil
}

func (nn nullNotifier) Send(msg Message, messageArgs Args) error {
	return nil
}

type argsMap struct {
	items map[string]string
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
