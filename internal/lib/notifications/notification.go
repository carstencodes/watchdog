package notifications

import (
	"errors"
	"flag"
	"fmt"

	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
)

type notificationVar struct {
	notifier string
}

func (n notificationVar) String() string {
	return n.notifier
}

func (n notificationVar) Set(value string) error {
	if _, ok := notificationClients[value]; ok {
		n.notifier = value
		return nil
	}

	return errors.New("Unknown notifier " + value)
}

var notificationClient = notificationVar{""}
var notificationClients = make(map[string]NotifierCreatorFunc)

func init() {
	flag.Var(&notificationClient, "notify-client", "The notification client to use")
}

func GetNotificationClient(log watchLog.Log) (Notifier, error) {
	creator, present := notificationClients[notificationClient.notifier]

	if !present {
		return nil, errors.New(fmt.Sprintf("Unknown notification client: %s", notificationClient))
	}

	return creator(log)
}
