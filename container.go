package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	docker_container "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type container struct {
	id       string
	name     string
	running  bool
	ignored  bool
	disabled bool
	healthy  bool
}

type containers struct {
	client    *client.Client
	ctx       *context.Context
	collector *collector
	logger    *log.Logger
	notifier  Notifier
	collect   collector
	items     []container
}

var syncRoot = sync.Mutex{}

const ignoreLabel = "com.github.carstencodes.watchtower.ignore"

func newContainersClient(col *collector, logger *log.Logger, notifier Notifier, coll collector, ctx *context.Context) (*containers, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	container_list := []container{}

	client := &containers{cli, ctx, col, logger, notifier, coll, container_list}

	return client, nil
}

func (cont *containers) updateContainers() error {
	list, err := cont.client.ContainerList(context.Background(), docker_container.ListOptions{
		All: true,
	})

	if err != nil {
		return err
	}

	count := len(list)

	syncRoot.Lock()
	cont.items = []container{}
	syncRoot.Unlock()

	items := make(chan container, count)
	wg := &sync.WaitGroup{}

	wg.Add(count)

	for _, cnt := range list {
		go cont.parse_container(items, wg, cnt)
	}

	wg.Wait()

	disabled := 0
	running := 0
	ignored := 0
	unhealthy := 0

	for item := range items {
		syncRoot.Lock()
		cont.items = append(cont.items, item)
		syncRoot.Unlock()
		if item.disabled {
			disabled += 1
		}
		if item.ignored {
			ignored += 1
		}
		if item.running {
			running += 1
		}
		if !item.healthy {
			unhealthy += 1
		}
	}

	cont.collect.metrics.disabled_containers.Set(float64(disabled))
	cont.collect.metrics.running_containers.Set(float64(running))
	cont.collect.metrics.ignored_containers.Set(float64(ignored))
	cont.collect.metrics.unhealthy_containers.Set(float64(unhealthy))

	return nil
}

func (cont *containers) parse_container(items chan<- container, wg *sync.WaitGroup, containerItem types.Container) {
	result := container{}

	defer wg.Done()

	defer func() {
		items <- result
	}()

	result.running = false
	result.ignored = true
	result.disabled = true
	result.healthy = true

	result.id = containerItem.ID
	if len(containerItem.Names) > 0 {
		result.name = containerItem.Names[0]
	} else {
		result.name = ""
	}

	disable, found := containerItem.Labels[ignoreLabel]
	result.disabled = found && disable == "false"

	inspect, err := cont.client.ContainerInspect(*cont.ctx, containerItem.ID)
	if err != nil {
		// TODO log
		return
	}

	state := inspect.State
	if state != nil {
		result.running = state.Running
		result.ignored = state.Health == nil
		result.healthy = state.Health.Status == types.Healthy || state.Health.Status == types.Starting
	}
}

func (cont *containers) refresh() {
	syncRoot.Lock()
	wg := &sync.WaitGroup{}
	wg.Add(len(cont.items))
	for _, containerItem := range cont.items {
		go cont.update_container(containerItem, wg)
	}
	syncRoot.Unlock()
	wg.Wait()
}

func (cont *containers) update_container(containerItem container, wg *sync.WaitGroup) {
	defer wg.Done()

	inspect, err := cont.client.ContainerInspect(*cont.ctx, containerItem.id)
	if err != nil {
		// TODO log
		return
	}

	state := inspect.State
	if state != nil {
		containerItem.running = state.Running
		containerItem.ignored = state.Health == nil
		containerItem.healthy = state.Health.Status == types.Healthy || state.Health.Status == types.Starting
	}
}

func (cont *containers) restartPending() {
	syncRoot.Lock()
	restartables := []container{}
	for _, containerItem := range cont.items {
		if (!containerItem.disabled || !containerItem.ignored) && containerItem.running && !containerItem.healthy {
			restartables = append(restartables, containerItem)
		}
	}
	syncRoot.Unlock()

	msg := "The following containers had to be restarted:\n\n"

	wg := &sync.WaitGroup{}
	wg.Add(len(restartables))

	for _, item := range restartables {
		go cont.restartContainer(item, wg)
		msg = msg + fmt.Sprintf("'%s' (%s)\n", item.name, item.id)
	}

	wg.Wait()

	args := argsMap{map[string]string{}}

	cont.notifier.Send(
		Message{
			title:   fmt.Sprintf("Watchdog restarted %d containers", len(restartables)),
			message: msg,
		},
		args,
	)
}

func (cont *containers) restartContainer(containerItem container, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := 2 * time.Minute
	t := int(timeout.Seconds())
	options := docker_container.StopOptions{
		Signal:  "",
		Timeout: &t,
	}
	err := cont.client.ContainerRestart(*cont.ctx, containerItem.id, options)
	if err != nil {
		// TODO log
	}

	if containerItem.name != "" {
		// TODO measure
	}

	cont.collect.metrics.restarted_containers.Inc()
}
