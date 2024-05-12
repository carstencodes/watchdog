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

package container

import (
	"context"

	"github.com/docker/docker/client"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	watchNotifier "github.com/carstencodes/watchdog/internal/lib/notifications"
)

type containerProxy struct {
	id       string
	name     string
	running  bool
	ignored  bool
	disabled bool
	healthy  bool
}

type ContainerCollection interface {
	UpdateContainers() error
	Refresh()
	RestartPending()
}

type containerCollectionImpl struct {
	client    *client.Client
	ctx       *context.Context
	collector watchCollector.Collector
	logger    watchLog.Log
	notifier  watchNotifier.Notifier
	items     []containerProxy
}
