package mockgen_first_run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/mockgen"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/errgroup"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type GenerateMocks struct {
	Fs afero.Fs

	DomainsPath string

	WiremockPath string

	Logger io.Writer
}

func NewMocksGenWithDefaultFs(domainsPath, wiremockPath string, logger io.Writer) GenerateMocks {
	return GenerateMocks{DomainsPath: domainsPath, WiremockPath: wiremockPath, Logger: logger, Fs: afero.NewOsFs()}
}

func (g *GenerateMocks) GenerateForEachDomain(ctx context.Context) error {
	domainPaths, err := fsutils.GatherDirs(g.Fs, g.DomainsPath)
	if err != nil {
		return fmt.Errorf("gather dirs: %w", err)
	}

	contractTypes := []types.SourceFileType{types.ProtoType, types.OpenAPIType}

	errG, errCtx := errgroup.WithContext(ctx)
	errG.SetLimit(len(domainPaths))

	for _, domainEntry := range domainPaths {
		domainName := domainEntry.Name()
		domainPath := filepath.Join(g.DomainsPath, domainName)

		should, err := g.shouldSkipDomain(domainName)
		if err != nil {
			return fmt.Errorf("should skip: %w", err)
		}

		if should {
			log.Printf("mocks generation: domain: '%s'. Skip\n", domainName)
			continue
		}

		log.Printf("mocks generation: domain: '%s'\n", domainName)

		generator := mockgen.NewMocksGenWithDefaultFs(domainPath, g.WiremockPath, g.Logger)

		for _, contractType := range contractTypes {
			if err := generator.Generate(errCtx, contractType, domainName); err != nil {
				log.Printf("mockgen: type '%s', domain '%s': %s. Skip\n", contractType, domainName, err)

				if err = handleMockgenErrors(err); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (g *GenerateMocks) shouldSkipDomain(domain string) (bool, error) {
	mocksPath := filepath.Join(g.WiremockPath, domain, "mappings")

	entries, err := afero.ReadDir(g.Fs, mocksPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("read dir: %w", err)
	}

	if len(entries) > 0 {
		return true, nil
	}

	return false, nil
}

func handleMockgenErrors(err error) error {
	if errors.Is(err, sourcer.ErrProvidedDirectoryHasNoContracts) {
		return nil
	}

	return fmt.Errorf("generaete: %w", err)
}
