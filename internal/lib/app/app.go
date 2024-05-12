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

package app

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
	"github.com/carstencodes/watchdog/internal/lib/common"
	watchContainer "github.com/carstencodes/watchdog/internal/lib/container"
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	watchNotifier "github.com/carstencodes/watchdog/internal/lib/notifications"
	watchWorker "github.com/carstencodes/watchdog/internal/lib/worker"
)

type App struct {
	context  appContext
	services appServices
	details  common.ApplicationDetails
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
	details := common.ApplicationInfo()
	err := setupConfiguration()
	if err != nil {
		return nil, err
	}

	var lg watchLog.Log
	lg, err = watchLog.CreateLog(details)
	if err != nil {
		return nil, err
	}
	lg.Info().Printf("Starting %s %s, %s", details.Name(), details.Version(), details.Copyright())

	var col = watchCollector.NewCollector(lg)
	var notifier watchNotifier.Notifier
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

	app := App{appCtx, appSvc, details}

	go func() {
		<-c
		app.Terminate()
	}()

	return &app, nil
}

func (app App) Run() error {
	flag.Parse()
	app.services.logger.Debug().Printf("Initialize collector.")
	app.services.collector.Init()
	app.services.logger.Debug().Printf("Initialize container collection.")
	err := app.services.containers.UpdateContainers()
	if err != nil {
		return err
	}
	app.services.logger.Debug().Printf("Refreshing container list.")
	app.services.containers.Refresh()
	app.services.logger.Debug().Printf("Starting worker.")
	app.services.worker.Start()
	app.services.logger.Debug().Printf("Starting Worker.")
	err = app.services.collector.Server().StartServer(app.context.ctx)
	if err != nil {
		return err
	}

	app.services.logger.Info().Printf("Application started successfully.")

	<-(*app.context.ctx).Done()
	app.services.logger.Debug().Printf("Stopping worker.")
	app.services.worker.Stop()
	app.services.logger.Info().Printf("Application terminated.")
	return nil
}

func (app App) Terminate() {
	app.services.logger.Info().Printf("Shutdown signal received. Terminating application.")
	app.context.cancel()
	app.services.logger.Debug().Printf("Context canceled.")
}
