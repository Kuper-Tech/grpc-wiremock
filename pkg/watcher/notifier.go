package watcher

import (
	"github.com/farmergreg/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

type Notifier interface {
	Errors() chan error
	Events() chan fsnotify.Event
}

type RecursiveNotifier struct {
	notifier *rfsnotify.RWatcher
}

func (e RecursiveNotifier) Errors() chan error {
	return e.notifier.Errors
}

func (e RecursiveNotifier) Events() chan fsnotify.Event {
	return e.notifier.Events
}

type NonRecursiveNotifier struct {
	notifier *fsnotify.Watcher
}

func (e NonRecursiveNotifier) Errors() chan error {
	return e.notifier.Errors
}

func (e NonRecursiveNotifier) Events() chan fsnotify.Event {
	return e.notifier.Events
}
