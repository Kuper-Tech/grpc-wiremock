package svcrunner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Run(ctx context.Context, path, port string) (func(), error) {
	process, err := start(ctx, path, port)
	if err != nil {
		return nil, fmt.Errorf("process start is failed: %w", err)
	}

	turnOffF := func() {
		log.Printf("Process %d and his whole family "+
			"are on their way to Valhalla", process.Pid)

		killGroup(process)
	}

	log.Printf("Process %d is started\n", process.Pid)

	return turnOffF, nil
}

func killGroup(process *os.Process) {
	if process == nil {
		log.Println("Process is nil. Can't stop the process")
		return
	}

	negativePIDToKillEntireGroup := -process.Pid

	// kill the entire group with child processes
	if err := syscall.Kill(negativePIDToKillEntireGroup, syscall.SIGKILL); err != nil {
		log.Printf("failed to kill child process %d\n", process.Pid)
	}
}

func start(ctx context.Context, path, port string) (*os.Process, error) {
	command := exec.CommandContext(ctx, "make", "run", "-C", path)
	command.Env = createEnvs(port)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("run cmd: %w", err)
	}

	return command.Process, nil
}

func createEnvs(port string) []string {
	const portEnv = "GRPC_TO_HTTP_PROXY_PORT"

	env := fmt.Sprintf("%s=%s", portEnv, port)
	envs := os.Environ()

	return append(envs, env)
}
