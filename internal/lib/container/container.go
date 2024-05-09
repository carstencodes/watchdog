package container

import (
	"context"
	"fmt"
	"sync"
	"time"

	watchCollector "github.com/carstencodes/watchdog/internal/lib/collector"
	watchLog "github.com/carstencodes/watchdog/internal/lib/log"
	watchNotifier "github.com/carstencodes/watchdog/internal/lib/notifications"

	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var syncRoot = sync.Mutex{}

const ignoreLabel = "com.github.carstencodes.watchtower.ignore"

func NewContainersClient(col watchCollector.Collector, logger watchLog.Log, notifier watchNotifier.Notifier, ctx *context.Context) (ContainerCollection, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	container_list := []containerProxy{}

	collection := &containerCollectionImpl{cli, ctx, col, logger, notifier, container_list}

	return collection, nil
}

func (cont containerCollectionImpl) UpdateContainers() error {
	list, err := cont.client.ContainerList(context.Background(), dockerContainer.ListOptions{
		All: true,
	})

	if err != nil {
		return err
	}

	count := len(list)

	syncRoot.Lock()
	cont.items = []containerProxy{}
	syncRoot.Unlock()

	items := make(chan containerProxy, count)
	wg := &sync.WaitGroup{}

	disabled := 0
	running := 0
	ignored := 0
	unhealthy := 0

	for _, cnt := range list {
		disable, found := cnt.Labels[ignoreLabel]
		isDisabled := found && disable == "false"
		if !isDisabled {
			wg.Add(1)
			go cont.parse_container(items, wg, cnt.ID)
		} else {
			disabled += 1
		}
	}

	wg.Wait()
	close(items)

	var resolvedItems []containerProxy
	for item := range items {
		resolvedItems = append(resolvedItems, item)
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

	syncRoot.Lock()
	cont.items = resolvedItems
	syncRoot.Unlock()

	cont.collector.CollectContainerStatistics(float64(disabled), float64(running), float64(ignored), float64(unhealthy))

	return nil
}

func (cont containerCollectionImpl) parse_container(items chan<- containerProxy, wg *sync.WaitGroup, containerId string) {
	result := containerProxy{}

	defer wg.Done()

	defer func() {
		items <- result
	}()

	result.running = false
	result.ignored = true
	result.disabled = true
	result.healthy = true

	result.id = containerId

	inspect, err := cont.client.ContainerInspect(*cont.ctx, containerId)
	if err != nil {
		// TODO log
		return
	}

	state := inspect.State
	if state != nil {
		result.running = state.Running
		result.ignored = state.Health == nil
		result.healthy = state.Health == nil || (state.Health.Status == types.Healthy || state.Health.Status == types.Starting)
	}

	if len(inspect.Name) > 0 {
		result.name = inspect.Name
	} else {
		result.name = ""
	}
}

func (cont containerCollectionImpl) Refresh() {
	syncRoot.Lock()
	wg := &sync.WaitGroup{}
	wg.Add(len(cont.items))
	for _, containerItem := range cont.items {
		go cont.update_container(containerItem, wg)
	}
	syncRoot.Unlock()
	wg.Wait()
}

func (cont containerCollectionImpl) update_container(containerItem containerProxy, wg *sync.WaitGroup) {
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

func (cont containerCollectionImpl) RestartPending() {
	syncRoot.Lock()
	var restartables []containerProxy
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

	args := watchNotifier.NewArgsMap(map[string]string{})

	_ = cont.notifier.Send(
		watchNotifier.NewMessage(
			fmt.Sprintf("Watchdog restarted %d containers", len(restartables)),
			msg,
		),
		args,
	)
}

func (cont containerCollectionImpl) restartContainer(containerItem containerProxy, wg *sync.WaitGroup) {
	defer wg.Done()
	timeout := 2 * time.Minute
	t := int(timeout.Seconds())
	options := dockerContainer.StopOptions{
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

	cont.collector.ContainerRestarted()
}
