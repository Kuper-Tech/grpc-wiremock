package proto

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/builder"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/builder/updaters"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/errgroup"
	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/loader"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/traverser"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/printer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type compiler interface {
	CompileToOpenAPI(context.Context, contract.Contract, string, io.Writer) error
}

type protoUnifier struct {
	fs afero.Fs

	logs io.Writer

	compiler
}

func NewUnifier(fs afero.Fs, compiler compiler, logs io.Writer) *protoUnifier {
	return &protoUnifier{fs: fs, compiler: compiler, logs: logs}
}

func (o *protoUnifier) Unify(ctx context.Context, contracts contract.SetOfContracts, output string) error {
	tmpPath := environment.TmpOverwrittenContractsDir

	tmpDirForContracts, err := afero.TempDir(o.fs, tmpPath, "")
	if err != nil {
		return fmt.Errorf("create tmp dir: %w", err)
	}

	goPackageUpdater := updaters.NewGoPackageUpdater()

	optionUpdater, err := updaters.NewOptionUpdater()
	if err != nil {
		return fmt.Errorf("create option updater: %w", err)
	}

	descriptors, err := traverser.Descriptors(contracts)
	if err != nil {
		return fmt.Errorf("get descriptors: %w", err)
	}

	updatedDescriptors, err := builder.UpdateContracts(descriptors, goPackageUpdater, optionUpdater)
	if err != nil {
		return fmt.Errorf("overwrite: %w", err)
	}

	if err = printer.Print(updatedDescriptors, tmpDirForContracts); err != nil {
		return fmt.Errorf("print: %w", err)
	}

	sourcerInstance, err := sourcer.New(o.fs, tmpDirForContracts, types.ProtoType)
	if err != nil {
		return fmt.Errorf("create sourcer for overwritten proto contracts: %w", err)
	}

	updatedContracts, err := loader.Load(o.fs, sourcerInstance)
	if err != nil {
		return fmt.Errorf("load overwritten proto contracts: %w", err)
	}

	errG, errCtx := errgroup.WithContext(ctx)
	errG.SetLimit(len(updatedContracts))

	for _, contractToCompile := range updatedContracts {
		contractToCompile := contractToCompile

		errG.Go(func() error {
			saveTo, err := o.createTmpDir(contractToCompile, output)
			if err != nil {
				return fmt.Errorf("create tmp dir: %w", err)
			}

			if err = o.CompileToOpenAPI(errCtx, contractToCompile, saveTo, o.logs); err != nil {
				return fmt.Errorf("compile: %w", err)
			}

			return nil
		})
	}

	if err = errG.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	if err = o.replaceFromTmpDirs(output); err != nil {
		return fmt.Errorf("replace: %w", err)
	}

	return nil
}

func (o *protoUnifier) createTmpDir(contract contract.Contract, output string) (string, error) {
	saveTo := filepath.Join(output, strings.TrimSuffix(filepath.Base(contract.HeaderPath), ".proto"))

	if err := o.fs.MkdirAll(saveTo, os.ModePerm); err != nil {
		return "", fmt.Errorf("mkdirall: %w", err)
	}

	return saveTo, nil
}

func (o *protoUnifier) replaceFromTmpDirs(output string) error {
	entries, err := fsutils.GatherEntries(o.fs, output)
	if err != nil {
		return fmt.Errorf("gather entries: %w", err)
	}

	var idx int
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "openapi.yaml" {
			newPath := filepath.Join(output, fmt.Sprintf("openapi%d.yaml", idx))
			if err := o.fs.Rename(entry.Path, newPath); err != nil {
				return fmt.Errorf("rename: %w", err)
			}
			idx++
		}
	}

	return nil
}
