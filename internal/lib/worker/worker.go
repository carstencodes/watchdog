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
