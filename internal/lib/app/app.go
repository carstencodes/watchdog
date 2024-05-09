package app

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
	watchContainer "github.com/carstencodes/watchdog/internal/lib/container"
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	watchNotifier "github.com/carstencodes/watchdog/internal/lib/notifications"
	watchWorker "github.com/carstencodes/watchdog/internal/lib/worker"
)

type App struct {
	ctx        *context.Context
	cancel     context.CancelFunc
	logger     *log.Logger
	collector  watchCollector.Collector
	notifier   watchNotifier.Notifier
	containers watchContainer.ContainerCollection
	worker     watchWorker.Worker
}

func NewApp() *App {
	var ctx = context.Background()
	var cancel context.CancelFunc
	var app *App
	var lg = watchLog.CreateLog()
	var col = watchCollector.NewCollector(lg)
	notifier, err := watchNotifier.GetNotificationClient()

	if err != nil {
		lg.Panicf("Failed to initialize notification client: %v", err)
	}

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt)

	containers, err := watchContainer.NewContainersClient(col, lg, notifier, &ctx)

	worker := watchWorker.CreateWorker(containers)

	app = &App{&ctx, cancel, lg, col, notifier, containers, worker}

	return app
}

func (app *App) Run() error {
	flag.Parse()
	app.collector.Init()
	err := app.containers.UpdateContainers()
	if err != nil {
		return err
	}
	app.containers.Refresh()

	app.worker.Start()

	err = app.collector.Server().StartServer(app.ctx)
	if err != nil {
		return err
	}

	<-(*app.ctx).Done()
	app.worker.Stop()
	return nil
}

func (app *App) Terminate() {
	app.cancel()
}
