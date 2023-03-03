package runner

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/utils/exec"
)

type shell struct {
	exec exec.Interface
}

func New(ex exec.Interface) shell {
	return shell{exec: ex}
}

func (s shell) Run(ctx context.Context, cmd string, args ...string) error {
	rawOut, err := s.exec.CommandContext(ctx, cmd, args...).CombinedOutput()
	if err == nil {
		return nil
	}

	out := strings.TrimSpace(string(rawOut))
	switch {
	case errors.Is(err, exec.ErrExecutableNotFound):
		return fmt.Errorf("'%s' not found on host: %w", cmd, err)
	default:
		return errUnknownCommandExecution{Command: cmd, Output: out, Args: args, Err: err}
	}
}
