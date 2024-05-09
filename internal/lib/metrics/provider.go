package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics interface {
	IncrementRestartedContainers()
	SetDisabledContainers(count float64)
	SetRunningContainers(count float64)
	SetUnhealthyContainers(count float64)
	SetIgnoredContainers(count float64)
}

type PrometheusSink interface {
	Register(registerer prometheus.Registerer)
}
