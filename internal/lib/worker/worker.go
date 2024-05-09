package worker

import (
	"time"

	"github.com/go-co-op/gocron"

	watchContainer "github.com/carstencodes/watchdog/internal/lib/container"
)

func CreateWorker(containers watchContainer.ContainerCollection) Worker {
	result := workerImpl{
		gocron.NewScheduler(time.Local),
		containers,
	}
	return result
}

func (w workerImpl) Start() {
	w.scheduler.Every(1).Hour().At(":00").Do(w.containers.UpdateContainers)
	w.scheduler.Every(10).Minutes().At(":00").Do(w.containers.Refresh)
	w.scheduler.Every(10).Minutes().At(":59").Do(w.containers.RestartPending)
	w.scheduler.StartAsync()
}

func (w workerImpl) Stop() {
	w.scheduler.Stop()
	w.scheduler.Clear()
}
