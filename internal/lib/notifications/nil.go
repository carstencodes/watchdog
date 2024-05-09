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

func (nn nullNotifier) Send(_ Message, _ Args) error {
	return nil
}

func (nn nullNotifier) CreateDefaultArgs() Args {
	return NewArgsMap(map[string]string{})
}

func init() {
	notificationClients[""] = newNullNotifier
}
