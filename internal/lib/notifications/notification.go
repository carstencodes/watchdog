package notifications

import (
	"errors"
	"flag"
	"fmt"

	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
)

var notificationClient string
var notificationClients = make(map[string]NotifierCreatorFunc)

func init() {
	flag.StringVar(&notificationClient, "notify", "", "The notification client")
}

func GetNotificationClient(log watchLog.Log) (Notifier, error) {
	creator, present := notificationClients[notificationClient]

	if !present {
		return nil, errors.New(fmt.Sprintf("Unknown notification client: %s", notificationClient))
	}

	return creator(log)
}
