package watcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/farmergreg/rfsnotify"
	"golang.org/x/time/rate"
	"gopkg.in/fsnotify.v1"
)

func NewRealWatcher(header WatcherDesc, logger io.Writer) (watcher, error) {
	var allowedNameRules []*regexp.Regexp

	for _, rule := range header.Behave.Entry.NameRules {
		compiledRule, err := regexp.Compile(rule)
		if err != nil {
			return watcher{}, fmt.Errorf("compile name rule %s: %w", rule, err)
		}

		allowedNameRules = append(allowedNameRules, compiledRule)
	}

	notifier, err := createNotifier(header.Recursive, header.Path)
	if err != nil {
		return watcher{}, fmt.Errorf("create notifier '%s': %w", header.Name, err)
	}

	limiter := createLimiter(header.Behave.Throttle)

	return watcher{
		WatcherDesc: header,

		notifier: notifier,
		limiter:  limiter,
		logger:   logger,

		allowedNameRules: allowedNameRules,
	}, nil
}

func (w *Watchers) Watch(ctx context.Context) error {
	if len(*w) == 0 {
		return fmt.Errorf("watchers must be provided")
	}

	var wg sync.WaitGroup

	for _, watcherToRun := range *w {
		watcherToRun := watcherToRun

		log.Printf("watcher '%s' is ready\n", watcherToRun.Name)

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := watcherToRun.watch(ctx); err != nil {
				log.Printf("watcher '%s': %s", watcherToRun.Name, err)
			}
		}()
	}

	wg.Wait()

	return nil
}

func createNotifier(recursive bool, path string) (Notifier, error) {
	if recursive {
		notifier, err := rfsnotify.NewWatcher()
		if err != nil {
			return nil, fmt.Errorf("create watchers: %w", err)
		}

		if err = notifier.AddRecursive(path); err != nil {
			return nil, fmt.Errorf("add recursive dir: %w", err)
		}

		return RecursiveNotifier{notifier}, nil
	}

	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create watchers: %w", err)
	}

	if err = notifier.Add(path); err != nil {
		return nil, fmt.Errorf("add dir: %w", err)
	}

	return NonRecursiveNotifier{notifier}, nil
}

func createLimiter(throttle ThrottlingRules) *rate.Limiter {
	const defaultBurst = 1

	return rate.NewLimiter(
		rate.Every(throttle.Interval),
		defaultBurst,
	)
}

func (w *watcher) watch(ctx context.Context) error {
	for {
		select {
		case event, ok := <-w.notifier.Events():
			if !ok {
				continue
			}

			if err := w.handleEvent(ctx, event); err != nil {
				if err = handleErrors(err); err != nil {
					return err
				}

				continue
			}

		case err, ok := <-w.notifier.Errors():
			if !ok {
				continue
			}

			if err != nil {
				return fmt.Errorf("notifier: %w", err)
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func (w *watcher) handleEvent(ctx context.Context, event fsnotify.Event) error {
	path := event.Name

	_, allowed := w.Behave.Event[event.Op]
	if !allowed {
		return skipEventErr
	}

	if err := w.filterNames(path); err != nil {
		return skipEventErr
	}

	if !w.limiter.Allow() {
		return skipEventErr
	}

	retryableFunc := func() error {
		if err := w.Do(ctx, w.Path, path); err != nil {
			return err
		}

		return nil
	}

	opts := []retry.Option{retry.Attempts(w.Behave.Retry.Attempts)}

	time.AfterFunc(w.Behave.Throttle.DelayAfterEvent, func() {
		if err := retry.Do(retryableFunc, opts...); err != nil {
			log.Printf("do with retry: %s", err)
		}
	})

	return nil
}

func (w *watcher) filterNames(path string) error {
	var atLeastOne bool

	for _, rule := range w.allowedNameRules {
		if rule.MatchString(path) {
			atLeastOne = true
			continue
		}
	}

	if !atLeastOne {
		return fmt.Errorf("name doesn't match")
	}

	return nil
}

func handleErrors(err error) error {
	if errors.Is(err, skipEventErr) {
		return nil
	}

	return err
}
