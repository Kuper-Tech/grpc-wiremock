package supervisord

import (
	"context"
	"fmt"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

type commandRunner interface {
	Run(ctx context.Context, cmd string, args ...string) error
}

type Reloader struct {
	Command []string
	Runner  commandRunner
}

func (r Reloader) ReloadConfig(ctx context.Context) error {
	if err := r.Runner.Run(ctx, sliceutils.FirstOf(r.Command), r.Command[1:]...); err != nil {
		return fmt.Errorf("reload: %w", err)
	}

	return nil
}
