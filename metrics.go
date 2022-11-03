package main

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	running_containers   prometheus.Gauge
	ignored_containers   prometheus.Gauge
	disabled_containers  prometheus.Gauge
	restarted_containers prometheus.Counter
	unhealthy_containers prometheus.Gauge
}

func newMetrics() *metrics {
	metrics := &metrics{
		running_containers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "running",
				Help:      "The amount of containers currently running on this machine",
			},
		),
		ignored_containers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "ignored",
				Help:      "The number of containers ignored due to a missing health-check",
			},
		),
		disabled_containers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "disabled",
				Help:      "The number of containers ignored due to a container or image label",
			},
		),
		restarted_containers: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "watchdog",
				Subsystem: "operations",
				Name:      "restarts",
				Help:      "The number of containers that were restarted by the engine",
			},
		),
		unhealthy_containers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "watchdog",
				Subsystem: "containers",
				Name:      "unhealthy",
				Help:      "The number of containers that are currently unhealthy",
			},
		),
	}

	return metrics
}

func (met *metrics) register(registerer prometheus.Registerer) {
	registerer.MustRegister(
		met.running_containers,
		met.disabled_containers,
		met.ignored_containers,
		met.unhealthy_containers,
		met.restarted_containers,
	)
}
