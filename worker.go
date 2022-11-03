package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

type worker struct {
	scheduler *gocron.Scheduler
}

func createWorker() *worker {
	result := &worker{
		gocron.NewScheduler(time.Local),
	}
	return result
}

func (w *worker) start(c *containers) {
	w.scheduler.Every(1).Hour().At(":00").Do(c.updateContainers)
	w.scheduler.Every(10).Minutes().At(":00").Do(c.refresh)
	w.scheduler.Every(10).Minutes().At(":59").Do(c.restartPending)
	w.scheduler.StartAsync()
}

func (w *worker) stop() {
	w.scheduler.Stop()
	w.scheduler.Clear()
}
