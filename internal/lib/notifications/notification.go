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
