package container

import (
	"context"
	"log"

	"github.com/docker/docker/client"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
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
	logger    *log.Logger
	notifier  watchNotifier.Notifier
	items     []containerProxy
}
