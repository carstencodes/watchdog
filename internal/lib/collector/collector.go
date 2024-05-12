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
