package mockgen

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/mocks"
	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/loader"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/unifier"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/runner"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

type GenerateMocks struct {
	Fs afero.Fs

	Input string

	Output string

	Logger io.Writer
}

func NewMocksGenWithDefaultFs(input, output string, logger io.Writer) GenerateMocks {
	return GenerateMocks{Input: input, Output: output, Logger: logger, Fs: afero.NewOsFs()}
}

func (g GenerateMocks) Generate(ctx context.Context, fileType types.SourceFileType, domain string) error {
	if err := environment.CleanTmpDirs(g.Fs); err != nil {
		return fmt.Errorf("clean tmp dirs: %w", err)
	}

	if err := environment.DumpProtos(g.Fs); err != nil {
		return fmt.Errorf("dump: %w", err)
	}

	unifiedContracts, err := g.unifyContracts(ctx, fileType)
	if err != nil {
		return fmt.Errorf("unify contracts: %w", err)
	}

	generator := mocks.NewGenerator(g.Fs, g.Output, g.Logger)

	if err = generator.Generate(ctx, unifiedContracts, domain); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}

func (g GenerateMocks) unifyContracts(ctx context.Context, fileType types.SourceFileType) (contract.SetOfContracts, error) {
	sourcerInstance, err := sourcer.New(g.Fs, g.Input, fileType)
	if err != nil {
		return nil, fmt.Errorf("create sourcer: %w", err)
	}

	contracts, err := loader.Load(g.Fs, sourcerInstance)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	commandRunner := runner.New(exec.New())

	contractsCompiler, err := compiler.New(g.Fs, commandRunner, environment.TmpAnnotationProtosDir)
	if err != nil {
		return nil, fmt.Errorf("create compiler: %w", err)
	}

	unifierInstance := unifier.NewUnifier(g.Fs, contractsCompiler, g.Logger)

	path := environment.TmpUnifiedContractsDir

	tmpDirForContract, err := afero.TempDir(g.Fs, path, "")
	if err != nil {
		return nil, fmt.Errorf("create tmp dir: %w", err)
	}

	if err = unifierInstance.Unify(ctx, contracts, tmpDirForContract); err != nil {
		return nil, fmt.Errorf("unify: %w", err)
	}

	sourcerInstance, err = sourcer.New(g.Fs, path, types.OpenAPIType)
	if err != nil {
		return nil, fmt.Errorf("create sourcer for unified contracts: %w", err)
	}

	contracts, err = loader.Load(g.Fs, sourcerInstance)
	if err != nil {
		return nil, fmt.Errorf("load unified contracts: %w", err)
	}

	return contracts, nil
}
