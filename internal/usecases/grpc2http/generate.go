package grpc2http

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/builder"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/builder/updaters"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/compiler"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/proxy"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract/loader"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract/traverser"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/printer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/renderer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/runner"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

type GenerateProxyUsecase struct {
	Path    string
	Output  string
	BaseURL string

	Logs io.Writer

	Fs afero.Fs
}

func NewProxyGen(path, output, baseURL string, logs io.Writer) GenerateProxyUsecase {
	return GenerateProxyUsecase{
		Path:    path,
		Output:  output,
		BaseURL: baseURL,
		Fs:      afero.NewOsFs(),
		Logs:    logs,
	}
}

func (p GenerateProxyUsecase) Generate(ctx context.Context) error {
	if err := environment.CleanTmpDirs(p.Fs); err != nil {
		return fmt.Errorf("clean tmp dirs: %w", err)
	}

	if err := environment.DumpProtos(p.Fs); err != nil {
		return fmt.Errorf("dump: %w", err)
	}

	commandRunner := runner.New(exec.New())

	contractsCompiler, err := compiler.New(p.Fs, commandRunner)
	if err != nil {
		return fmt.Errorf("create compiler: %w", err)
	}

	const pathToTemplates = "proxy/files"

	projectRenderer, err := renderer.New(static.FromEmbed(), pathToTemplates)
	if err != nil {
		return fmt.Errorf("create renderer: %w", err)
	}

	contracts, err := p.overwriteContracts()
	if err != nil {
		return fmt.Errorf("prepare contracts: %w", err)
	}

	generator, err := proxy.NewGenerator(
		p.Fs, contractsCompiler,
		projectRenderer, p.BaseURL, p.Output, p.Logs,
	)
	if err != nil {
		return fmt.Errorf("create proxy generator: %w", err)
	}

	if err = generator.Generate(ctx, contracts); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err = p.cleanUp(); err != nil {
		return fmt.Errorf("clean up: %w", err)
	}

	return nil
}

func (p GenerateProxyUsecase) overwriteContracts() (protocontract.SetOfContracts, error) {
	contractsSourcer, err := sourcer.New(p.Fs, p.Path, types.ProtoType)
	if err != nil {
		return nil, fmt.Errorf("create sourcer: %w", err)
	}

	contracts, err := loader.Load(p.Fs, contractsSourcer)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	goPackageUpdater := updaters.NewGoPackageUpdater()

	descriptors := traverser.Descriptors(contracts)

	updatedDescriptors, err := builder.UpdateContracts(descriptors, goPackageUpdater)
	if err != nil {
		return nil, fmt.Errorf("overwrite: %w", err)
	}

	path := environment.TmpOverwrittenContractsDir

	tmpDirForContract, err := afero.TempDir(p.Fs, path, "")
	if err != nil {
		return nil, fmt.Errorf("create tmp dir: %w", err)
	}

	if err = printer.Print(updatedDescriptors, tmpDirForContract); err != nil {
		return nil, fmt.Errorf("print: %w", err)
	}

	contractsSourcer, err = sourcer.New(p.Fs, tmpDirForContract, types.ProtoType)
	if err != nil {
		return nil, fmt.Errorf("create sourcer for overwritten contracts: %w", err)
	}

	contracts, err = loader.Load(p.Fs, contractsSourcer)
	if err != nil {
		return nil, fmt.Errorf("load overwritten contracts: %w", err)
	}

	return contracts, nil
}

func (p GenerateProxyUsecase) cleanUp() error {
	if err := p.removeTmpFilesFromProxy(); err != nil {
		return fmt.Errorf("remove .keep files: %w", err)
	}

	if err := p.renameGoMod(); err != nil {
		return fmt.Errorf("rename go.mod: %w", err)
	}

	return nil
}

func (p GenerateProxyUsecase) removeTmpFilesFromProxy() error {
	const keepFile = ".keep"
	return fsutils.RemoveWithSubdirs(p.Fs, p.Output, keepFile)
}

func (p GenerateProxyUsecase) renameGoMod() error {
	const (
		goModName    = "go.mod"
		tmpGoModName = "go.mod.rename.me"
	)

	var (
		oldPath = filepath.Join(p.Output, tmpGoModName)
		newPath = filepath.Join(p.Output, goModName)
	)

	if err := p.Fs.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}
