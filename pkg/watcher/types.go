package watcher

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/fsnotify.v1"
)

type Action func(context.Context, string, string) error

type Watchers []watcher

var skipEventErr = fmt.Errorf("skip event")

type WatcherDesc struct {
	Do Action

	Name string
	Path string

	Behave BehaviourDesc

	Recursive bool
}

type watcher struct {
	WatcherDesc

	notifier Notifier

	limiter *rate.Limiter

	allowedNameRules []*regexp.Regexp

	logger io.Writer
}

type EventTypes map[fsnotify.Op]struct{}

func NewEventTypes() EventTypes {
	return map[fsnotify.Op]struct{}{}
}

func (e EventTypes) WithCreate() EventTypes {
	e[fsnotify.Create] = struct{}{}
	return e
}

func (e EventTypes) WithRemove() EventTypes {
	e[fsnotify.Remove] = struct{}{}
	return e
}

func (e EventTypes) WithRename() EventTypes {
	e[fsnotify.Rename] = struct{}{}
	return e
}

func (e EventTypes) WithWrite() EventTypes {
	e[fsnotify.Write] = struct{}{}
	return e
}

func (e EventTypes) WithChmod() EventTypes {
	e[fsnotify.Chmod] = struct{}{}
	return e
}

type EntryNamesRules []string

type EntryRules struct {
	NameRules EntryNamesRules
}

func NewEntryRules() EntryRules {
	return EntryRules{}
}

func (e EntryRules) WithNameRule(rule string) EntryRules {
	e.NameRules = append(e.NameRules, rule)
	return e
}

type ThrottlingRules struct {
	DelayAfterEvent time.Duration
	Interval        time.Duration
}

type RetryRules struct {
	Attempts uint
}

type BehaviourDesc struct {
	Event EventTypes
	Entry EntryRules
	Retry RetryRules

	Throttle ThrottlingRules
}
