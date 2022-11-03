package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort int    = 8080
	endpoint    string = "/metrics"
)

var (
	port int
)

var (
	protect_server sync.Mutex = sync.Mutex{}
)

func init() {
	flag.IntVar(&port, "port", defaultPort, "The tcp port to listen to")
}

type collector struct {
	logger  *log.Logger
	metrics *metrics
	reg     *prometheus.Registry
	server  *http.Server
}

func NewCollector(logger *log.Logger) *collector {
	col := &collector{
		logger,
		newMetrics(),
		prometheus.NewRegistry(),
		nil,
	}

	return col
}

func (col *collector) init() {
	col.metrics.register(col.reg)
	col.reg.MustRegister(collectors.NewBuildInfoCollector())
	col.reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(
			collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))
}

func (col *collector) startServer(ctx *context.Context) error {
	const listenerAddressFamily string = "tcp"

	addr := fmt.Sprintf(":%d", port)
	httpHandler := promhttp.HandlerFor(
		col.reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)
	mux := http.NewServeMux()

	listener, err := net.Listen(listenerAddressFamily, addr)
	if err != nil {
		return err
	}

	mux.Handle(endpoint, httpHandler)

	protect_server.Lock()
	col.server = &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: col.logger,
		BaseContext: func(_ net.Listener) context.Context {
			return *ctx
		},
	}
	protect_server.Unlock()

	go runServer(col.server, listener)

	return nil
}

func (col *collector) stopServer(ctx *context.Context) error {
	if col.server != nil {
		protect_server.Lock()
		err := col.server.Shutdown(*ctx)
		col.server = nil
		protect_server.Unlock()
		return err
	}

	return nil
}

func runServer(server *http.Server, listener net.Listener) error {
	server_error := server.Serve(listener)
	if !errors.Is(server_error, http.ErrServerClosed) {
		return server_error
	}

	return nil
}
