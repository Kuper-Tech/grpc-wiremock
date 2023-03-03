package executils

import (
	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

func MakeFakeCmd(fakeCmd *testingexec.FakeCmd, cmd string, args ...string) testingexec.FakeCommandAction {
	c := cmd
	a := args
	return func(cmd string, args ...string) exec.Cmd {
		return testingexec.InitFakeCmd(fakeCmd, c, a...)
	}
}

func MakeFakeOutput(output string, err error) testingexec.FakeAction {
	o := output
	return func() ([]byte, []byte, error) {
		return []byte(o), nil, err
	}
}
