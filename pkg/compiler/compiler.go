package compiler

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler/command"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler/compilecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/executils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type runner interface {
	Run(ctx context.Context, cmd string, args ...string) error
}

type compiler struct {
	compilerPath    string
	extraProtoPaths []string

	runner
}

func New(fs afero.Fs, shell runner, protoPaths ...string) (compiler, error) {
	protocPath, err := exec.LookPath(ProtoCompiler)
	if err != nil {
		return compiler{}, fmt.Errorf("look path for '%s': %w", ProtoCompiler, err)
	}

	if err = executils.HostHasBinaries(ProtoGoPlugin, ProtoGRPCPlugin, ProtoOpenAPIPlugin); err != nil {
		return compiler{}, fmt.Errorf("check binary: %w", err)
	}

	if err = fsutils.ValidDirectories(fs, protoPaths...); err != nil {
		return compiler{}, fmt.Errorf("check proto path: %w", err)
	}

	extraProtoPaths := []string{environment.TmpWellKnownProtosDir}
	extraProtoPaths = append(extraProtoPaths, protoPaths...)

	return compiler{
		runner:          shell,
		compilerPath:    protocPath,
		extraProtoPaths: extraProtoPaths,
	}, nil
}

func fromProtoContract(contract protocontract.Contract) compilecontract.Contract {
	return compilecontract.Contract{
		HeaderPath:   contract.HeaderPath,
		ImportsPaths: contract.ImportsPaths,
	}
}

func fromBaseContract(contract basecontract.Contract) compilecontract.Contract {
	return compilecontract.Contract{
		HeaderPath:   contract.HeaderPath,
		ImportsPaths: contract.ImportsPaths,
	}
}

func (c compiler) CompileToGo(ctx context.Context, contract protocontract.Contract, output string, logs io.Writer) error {
	commandToCompile := command.CompileGoPackages(output)
	return c.compile(ctx, fromProtoContract(contract), logs, commandToCompile)
}

func (c compiler) CompileToOpenAPI(ctx context.Context, contract basecontract.Contract, output string, logs io.Writer) error {
	commandToCompile := command.CompileOpenAPI(output)
	return c.compile(ctx, fromBaseContract(contract), logs, commandToCompile)
}

func (c compiler) compile(ctx context.Context, contract compilecontract.Contract, logs io.Writer, compileCmd command.CompileCommand) error {
	compile := command.BuildArgs(contract, c.extraProtoPaths, compileCmd)

	if _, err := fmt.Fprintln(logs, logProtocCall(c.compilerPath, compile)); err != nil {
		return fmt.Errorf("print: %w", err)
	}

	if err := c.Run(ctx, c.compilerPath, compile.Args()...); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func logProtocCall(protoc string, cmd command.Command) string {
	return fmt.Sprintf("%s \\%s", protoc, cmd.String())
}
