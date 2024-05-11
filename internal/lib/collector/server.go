package collector

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"

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
	protectServer = sync.Mutex{}
)

func init() {
	flag.IntVar(&port, "server-port", defaultPort, "The tcp port to listen to")
}

func (col collectorImpl) StartServer(ctx *context.Context) error {
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

	childContext, cancelWithCause := context.WithCancelCause(*ctx)

	protectServer.Lock()
	col.server = &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: col.logger.Error().GetLog(),
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(childContext, "watchdog.server.listener", listener)
		},
	}
	protectServer.Unlock()

	go runServer(col.server, listener, cancelWithCause)

	return nil
}

func (col collectorImpl) StopServer(ctx *context.Context) error {
	if col.server != nil {
		protectServer.Lock()
		err := col.server.Shutdown(*ctx)
		col.server = nil
		protectServer.Unlock()
		return err
	}

	return nil
}

func runServer(server *http.Server, listener net.Listener, causeFunc context.CancelCauseFunc) {
	serverError := server.Serve(listener)
	if !errors.Is(serverError, http.ErrServerClosed) {
		causeFunc(serverError)
	}
}
