package collector

import (
	"context"
)

type Server interface {
	StartServer(ctx *context.Context) error
	StopServer(ctx *context.Context) error
}

type Collector interface {
	Init()
	Server() Server
	CollectContainerStatistics(
		disabled float64,
		running float64,
		ignored float64,
		unhealthy float64)
	ContainerRestarted()
}
