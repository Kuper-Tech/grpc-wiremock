package loader

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

// contractLoader provides specific loader. For proto and OpenAPI contracts.
type contractLoader interface {
	LoadDir() (protocontract.SetOfContracts, error)
	LoadFile() (protocontract.SetOfContracts, error)
}

// contractSourcer provides validated source for requested contract/dir with contracts.
type contractSourcer interface {
	SourceInputs() []string
	SourceLoadType() types.SourceLoadType
	SourceFileType() types.SourceFileType
}

var unsupportedLoaderErr = fmt.Errorf("loader is not supported")

// Load allows to load contracts from provided source as `SetOfContracts`.
func Load(fs afero.Fs, source contractSourcer) (contracts protocontract.SetOfContracts, err error) {
	sourceType := source.SourceLoadType()

	specificLoader, err := getLoader(fs, source)
	if err != nil {
		return nil, fmt.Errorf("get loader: %w", err)
	}

	switch sourceType {
	case types.SourceDirType:
		contracts, err = specificLoader.LoadDir()
		if err != nil {
			return nil, fmt.Errorf("load contracts from dir: %w", err)
		}

	case types.SourceSingleType:
		contracts, err = specificLoader.LoadFile()
		if err != nil {
			return nil, fmt.Errorf("load contracts from file: %w", err)
		}

	default:
		return nil, fmt.Errorf("incorrect source type: %v", sourceType)
	}

	return contracts, nil
}

func getLoader(fs afero.Fs, s contractSourcer) (contractLoader, error) {
	switch s.SourceFileType() {
	case types.ProtoType:
		return NewLoader(fs, s), nil
	default:
		return nil, unsupportedLoaderErr
	}
}
