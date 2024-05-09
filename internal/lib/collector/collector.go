package collector

import (
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	watchMetrics "github.com/carstencodes/watchdog/internal/lib/metrics"
)

type collectorImpl struct {
	logger  watchLog.Log
	metrics watchMetrics.Metrics
	reg     *prometheus.Registry
	server  *http.Server
}

func NewCollector(logger watchLog.Log) Collector {
	col := &collectorImpl{
		logger,
		watchMetrics.NewMetrics(),
		prometheus.NewRegistry(),
		nil,
	}

	return col
}

func (col collectorImpl) Init() {
	if sink, ok := col.metrics.(watchMetrics.PrometheusSink); ok {
		sink.Register(col.reg)
	} else {
		// TODO Error logging
	}

	col.reg.MustRegister(collectors.NewBuildInfoCollector())
	col.reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))
}

func (col collectorImpl) Server() Server {
	return col
}

func (col collectorImpl) CollectContainerStatistics(
	disabled float64,
	running float64,
	ignored float64,
	unhealthy float64) {
	col.metrics.SetDisabledContainers(disabled)
	col.metrics.SetRunningContainers(running)
	col.metrics.SetIgnoredContainers(ignored)
	col.metrics.SetUnhealthyContainers(unhealthy)
}
func (col collectorImpl) ContainerRestarted() {
	col.metrics.IncrementRestartedContainers()
}
