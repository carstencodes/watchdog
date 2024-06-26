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
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"

	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
)

type gotifyConfig struct {
	baseUrl         string
	appToken        string
	defaultPriority int64
}

type gotifyClient struct {
	client *client.GotifyREST
	auth   *runtime.ClientAuthInfoWriter
	config gotifyConfig
}

type gotifyMessageArgs struct {
	priority int64
}

const defaultMessagePriority int64 = 5

var (
	clientConfig gotifyConfig
)

func init() {
	flag.StringVar(&clientConfig.appToken, "gotify-app-token", "", "The gotify app token to use, if notification client is set to 'gotify'")
	flag.StringVar(&clientConfig.baseUrl, "gotify-base-url", "", "The gotify base URL to send messages to, if notification client is set to 'gotify'")
	flag.Int64Var(&clientConfig.defaultPriority, "gotify-default-priority", 5, "The default message priority to use if notification client is set to 'gotify'")

	notificationClients["gotify"] = newGotifyClient
}

func newGotifyClient(_ watchLog.Log) (Notifier, error) {
	return fromGotifyClientConfig(clientConfig)
}

func fromGotifyClientConfig(config gotifyConfig) (Notifier, error) {
	if len(config.appToken) == 0 {
		return nil, errors.New("No app-token specified")
	}
	if len(config.baseUrl) == 0 {
		return nil, errors.New("No basic url specified")
	}

	baseUrl, parseErr := url.Parse(config.baseUrl)
	if parseErr != nil {
		return nil, parseErr
	}

	c := gotify.NewClient(baseUrl, &http.Client{})
	a := auth.TokenAuth(config.appToken)

	gotifySender := gotifyClient{c, &a, config}
	return gotifySender, nil
}

func (gotify gotifyClient) Connect() error {
	var limit int64 = 1

	_, err := gotify.client.Version.GetVersion(nil)
	if err != nil {
		return err
	}

	opts := message.NewGetMessagesParams()
	opts.Limit = &limit
	_, err = gotify.client.Message.GetMessages(
		opts,
		*gotify.auth,
	)

	return err
}

func (gotify gotifyClient) Disconnect() error {
	return nil
}

func (gotify gotifyClient) CreateDefaultArgs() Args {
	return newGotifyArgs(gotify.config.defaultPriority)
}

func (gotify gotifyClient) Send(msg Message, messageArgs Args) error {
	priority := getMessagePriority(messageArgs)

	messageParams := message.NewCreateMessageParams()
	messageParams.Body = &models.MessageExternal{
		Title:    msg.getTitle(),
		Message:  msg.getMessage(),
		Priority: int(priority),
	}

	_, messageErr := gotify.client.Message.CreateMessage(messageParams, *gotify.auth)

	return messageErr
}

func getMessagePriority(messageArgs Args) int64 {
	switch messageArgs.(type) {
	case gotifyMessageArgs:
		args := messageArgs.(gotifyMessageArgs)
		return args.priority
	default:
		setPrio, err := messageArgs.GetInt("priority")
		if err != nil {
			return defaultMessagePriority
		}
		return setPrio
	}
}

func newGotifyArgs(priority int64) Args {
	args := gotifyMessageArgs{priority}
	return args
}

func (args gotifyMessageArgs) GetBool(key string) (bool, error) {
	return false, errors.New(fmt.Sprintf("Failed to fetch key %s", key))
}

func (args gotifyMessageArgs) GetFloat(key string) (float64, error) {
	return 0.0, errors.New(fmt.Sprintf("Failed to fetch key %s", key))
}

func (args gotifyMessageArgs) GetInt(key string) (int64, error) {
	if key == "priority" {
		return args.priority, nil
	}

	return 0, errors.New(fmt.Sprintf("Failed to fetch key %s", key))
}

func (args gotifyMessageArgs) GetString(key string) (string, error) {
	return "", errors.New(fmt.Sprintf("Failed to fetch key %s", key))
}
