package notifications

import (
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
)

type nullNotifier struct{}

func newNullNotifier(_ watchLog.Log) (Notifier, error) {
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

func init() {
	notificationClients[""] = newNullNotifier
}
