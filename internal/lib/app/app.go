package app

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
	watchContainer "github.com/carstencodes/watchdog/internal/lib/container"
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	"github.com/carstencodes/watchdog/internal/lib/log/sinks"
	watchNotifier "github.com/carstencodes/watchdog/internal/lib/notifications"
	watchWorker "github.com/carstencodes/watchdog/internal/lib/worker"
)

type App struct {
	context  appContext
	services appServices
}

type appContext struct {
	ctx    *context.Context
	cancel context.CancelFunc
}

type appServices struct {
	logger     watchLog.Log
	notifier   watchNotifier.Notifier
	collector  watchCollector.Collector
	containers watchContainer.ContainerCollection
	worker     watchWorker.Worker
}

func NewApp() (*App, error) {
	var notifier watchNotifier.Notifier

	lg, err := watchLog.CreateLog(watchLog.Info, watchLog.NewSetup().WithSink(sinks.StdOut()))
	if err != nil {
		return nil, err
	}

	var col = watchCollector.NewCollector(lg)
	notifier, err = watchNotifier.GetNotificationClient(lg)

	if err != nil {
		lg.Fatal().Printf("Failed to initialize notification client: %v", err)
		return nil, err
	}

	var ctx = context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	containers, err := watchContainer.NewContainersClient(col, lg, notifier, &ctx)

	worker := watchWorker.CreateWorker(containers)

	appCtx := appContext{&ctx, cancel}
	appSvc := appServices{
		lg, notifier, col, containers, worker,
	}

	app := App{appCtx, appSvc}

	go func() {
		<-c
		app.Terminate()
	}()

	return &app, nil
}

func (app App) Run() error {
	flag.Parse()
	app.services.collector.Init()
	err := app.services.containers.UpdateContainers()
	if err != nil {
		return err
	}
	app.services.containers.Refresh()
	app.services.worker.Start()
	err = app.services.collector.Server().StartServer(app.context.ctx)
	if err != nil {
		return err
	}

	<-(*app.context.ctx).Done()
	app.services.worker.Stop()
	return nil
}

func (app App) Terminate() {
	app.context.cancel()
}