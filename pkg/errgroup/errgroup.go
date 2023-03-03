package errgroup

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Group struct {
	*errgroup.Group
}

func WithContext(parent context.Context) (*Group, context.Context) {
	g, ctx := errgroup.WithContext(parent)
	return &Group{g}, ctx
}

func (g *Group) Go(originFn func() error) {
	fn := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("recovered panic from gorutine: %v", r)
			}
		}()

		err = originFn()

		return
	}

	g.Group.Go(fn)
}
