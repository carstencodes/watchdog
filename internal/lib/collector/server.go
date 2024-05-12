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
