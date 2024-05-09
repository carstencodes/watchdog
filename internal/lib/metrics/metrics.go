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
