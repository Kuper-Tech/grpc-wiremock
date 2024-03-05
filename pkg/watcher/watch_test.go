package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

var (
	osfs = afero.NewOsFs()
	mu   sync.Mutex

	handledEvents []string
)

func resetHandledEvents() {
	mu.Lock()
	defer mu.Unlock()

	handledEvents = []string{}
}

func defaultGreeter(_ context.Context, _, path string) error {
	mu.Lock()
	defer mu.Unlock()

	handledEvents = append(handledEvents, path)

	// do something useful.

	return nil
}

const jsonFileRule = `.*\.json`

func TestWatchers_Watch(t *testing.T) {
	tests := []struct {
		name  string
		count int

		interval time.Duration

		watcherDesc WatcherDesc

		wantHandledEvents []string
	}{
		{
			watcherDesc: WatcherDesc{
				Do:   defaultGreeter,
				Name: "greeter-with-throttling",
				Path: "/tmp/test/watcher",
				Behave: BehaviourDesc{
					Event: NewEventTypes().WithCreate(),
					Entry: NewEntryRules().
						WithNameRule(jsonFileRule),
					Throttle: ThrottlingRules{Interval: 400 * time.Millisecond},
				},
			},

			count:    5,
			interval: time.Millisecond * 50,

			wantHandledEvents: []string{
				filepath.Join("/tmp/test/watcher", createName(0)),
			},
		},
		{
			watcherDesc: WatcherDesc{
				Do:   defaultGreeter,
				Name: "greeter-without-throttling",
				Path: "/tmp/test/watcher",
				Behave: BehaviourDesc{
					Event: NewEventTypes().WithCreate(),
					Entry: NewEntryRules().
						WithNameRule(jsonFileRule),
				},
			},

			count:    5,
			interval: time.Millisecond * 100,

			wantHandledEvents: createWantBuffer(5, "/tmp/test/watcher"),
		},
		{
			watcherDesc: WatcherDesc{
				Do:   defaultGreeter,
				Name: "greeter-only-rename-events",
				Path: "/tmp/test/watcher",
				Behave: BehaviourDesc{
					Event: NewEventTypes().WithRename(),
				},
			},

			count:    5,
			interval: time.Millisecond * 100,

			wantHandledEvents: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetHandledEvents()

			err := prepareEnvironment(tt.watcherDesc.Path)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())

			go createFiles(ctx, cancel, tt.watcherDesc.Path, tt.count, tt.interval)

			realWatcher, err := NewRealWatcher(tt.watcherDesc, os.Stdout)
			require.NoError(t, err)

			err = realWatcher.watch(ctx)
			require.NoError(t, err)

			require.ElementsMatch(t, handledEvents, tt.wantHandledEvents)
		})
	}
}

func prepareEnvironment(path string) error {
	if err := fsutils.RemoveTmpDirs(osfs, path); err != nil {
		return fmt.Errorf("remove tmp dirs: %w", err)
	}

	return nil
}

func createFiles(ctx context.Context, cancelFn func(), path string, count int, interval time.Duration) {
	defer cancelFn()

	var counter int

	ticker := time.Tick(interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("ctx.done")
			return
		case <-ticker:
			targetPath := filepath.Join(path, createName(counter))

			if err := afero.WriteFile(osfs, targetPath, []byte{}, os.ModePerm); err != nil {
				log.Printf("write file err: %s\n", err)
				return
			}

			if counter >= count {
				log.Println("counter.done")
				return
			}

			counter++
		}
	}
}

func createWantBuffer(count int, path string) []string {
	var wantBuffer []string

	for i := 0; i < count; i++ {
		wantBuffer = append(wantBuffer, filepath.Join(path, createName(i)))
	}

	return wantBuffer
}

func createName(idx int) string {
	return fmt.Sprintf("file_%d.json", idx)
}
