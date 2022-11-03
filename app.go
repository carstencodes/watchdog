package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
)

type App struct {
	ctx        *context.Context
	cancel     context.CancelFunc
	logger     *log.Logger
	collector  *collector
	notifier   *Notifier
	containers *containers
	worker     *worker
}

func NewApp() *App {
	var ctx context.Context = context.Background()
	var cancel context.CancelFunc
	var app *App
	var lg *log.Logger = createLog()
	var col *collector = NewCollector(lg)
	notifier, err := getNotificationClient()

	if err != nil {
		lg.Panicf("Failed to initialized notification client: %v", err)
	}

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt)

	containers, err := newContainersClient(col, lg, notifier, *col, &ctx)

	worker := createWorker()

	app = &App{&ctx, cancel, lg, col, &notifier, containers, worker}

	return app
}

func (app *App) Run() error {
	flag.Parse()
	app.collector.init()
	app.containers.updateContainers()
	app.containers.refresh()

	app.worker.start(app.containers)

	err := app.collector.startServer(app.ctx)
	if err != nil {
		return err
	}

	<-(*app.ctx).Done()
	app.worker.stop()
	return nil
}

func (app *App) Terminate() {
	app.cancel()
}
