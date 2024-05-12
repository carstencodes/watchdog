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

package metrics

import "github.com/prometheus/client_golang/prometheus"

type metricsImpl struct {
	runningContainers   prometheus.Gauge
	ignoredContainers   prometheus.Gauge
	disabledContainers  prometheus.Gauge
	restartedContainers prometheus.Counter
	unhealthyContainers prometheus.Gauge
}

func NewMetrics() Metrics {
	metrics := metricsImpl{
		runningContainers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "running",
				Help:      "The amount of containers currently running on this machine",
			},
		),
		ignoredContainers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "ignored",
				Help:      "The number of containers ignored due to a missing health-check",
			},
		),
		disabledContainers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "disabled",
				Help:      "The number of containers ignored due to a container or image label",
			},
		),
		restartedContainers: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "watchdog",
				Subsystem: "operations",
				Name:      "restarts",
				Help:      "The number of containers that were restarted by the engine",
			},
		),
		unhealthyContainers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "unhealthy",
				Help:      "The number of containers that are currently unhealthy",
			},
		),
	}

	return &metrics
}

func (m metricsImpl) Register(registerer prometheus.Registerer) {
	registerer.MustRegister(
		m.runningContainers,
		m.disabledContainers,
		m.ignoredContainers,
		m.unhealthyContainers,
		m.restartedContainers,
	)
}
