package watcher

import (
	"context"
	"fmt"
	"io"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/watcher"
)

type watchersRunner struct {
	logger io.Writer
}

func NewRunner(logger io.Writer) watchersRunner {
	return watchersRunner{logger: logger}
}

func (r *watchersRunner) Watch(ctx context.Context, requests ...WatchRequest) error {
	watchers, err := r.createWatchers(requests)
	if err != nil {
		return fmt.Errorf("create watchers: %w", err)
	}

	if err = watchers.Watch(ctx); err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	return nil
}

func (r *watchersRunner) createWatchers(requests []WatchRequest) (watcher.Watchers, error) {
	var watchers watcher.Watchers

	for _, request := range requests {
		preDefinedWatcher, exists := knownWatchers[request.Name]
		if !exists {
			return nil, fmt.Errorf("watcher '%s' is not found", request.Name)
		}

		preDefinedWatcher.Path = request.Path

		realWatcher, err := watcher.NewRealWatcher(preDefinedWatcher, r.logger)
		if err != nil {
			return nil, fmt.Errorf("create watcher: %w", err)
		}

		watchers = append(watchers, realWatcher)
	}

	return watchers, nil
}
