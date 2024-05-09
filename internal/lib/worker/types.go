package worker

import (
	"github.com/go-co-op/gocron"

	watchContainer "github.com/carstencodes/watchdog/internal/lib/container"
)

type workerImpl struct {
	scheduler  *gocron.Scheduler
	containers watchContainer.ContainerCollection
}

type Worker interface {
	Start()
	Stop()
}
