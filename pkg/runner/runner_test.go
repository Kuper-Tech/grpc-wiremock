package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/executils"
)

type execArgs struct {
	command string
	args    []string
	output  string
	err     error
}

func Test_shell_Run(t *testing.T) {
	tests := []struct {
		name string

		execScript execArgs

		want    string
		wantErr error
	}{
		{
			execScript: execArgs{"protoc", []string{"-I", "/tmp/protos", "baugi.proto"}, "", exec.ErrExecutableNotFound},
			wantErr:    exec.ErrExecutableNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExec := &testingexec.FakeExec{ExactOrder: true}

			fakeCmd := &testingexec.FakeCmd{}

			cmdAction := executils.MakeFakeCmd(fakeCmd, tt.execScript.command, tt.execScript.args...)
			outputAction := executils.MakeFakeOutput(tt.execScript.output, tt.execScript.err)

			fakeCmd.CombinedOutputScript = append(fakeCmd.CombinedOutputScript, outputAction)
			fakeExec.CommandScript = append(fakeExec.CommandScript, cmdAction)

			ctx := context.Background()
			err := New(fakeExec).Run(ctx, tt.execScript.command, tt.execScript.args...)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
