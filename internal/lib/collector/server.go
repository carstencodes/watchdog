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
	flag.IntVar(&port, "port", defaultPort, "The tcp port to listen to")
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

	protectServer.Lock()
	col.server = &http.Server{
		Addr:     addr,
		Handler:  mux,
		ErrorLog: col.logger,
		BaseContext: func(_ net.Listener) context.Context {
			return *ctx
		},
	}
	protectServer.Unlock()

	go runServer(col.server, listener)

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

func runServer(server *http.Server, listener net.Listener) error {
	server_error := server.Serve(listener)
	if !errors.Is(server_error, http.ErrServerClosed) {
		return server_error
	}

	return nil
}
