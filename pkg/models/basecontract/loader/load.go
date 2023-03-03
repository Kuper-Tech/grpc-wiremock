package loader

import (
	"fmt"

	"github.com/spf13/afero"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/loader/openapi"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/loader/proto"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

// contractLoader provides specific loader. For proto and OpenAPI contracts.
type contractLoader interface {
	LoadDir() (contract.SetOfContracts, error)
	LoadFile() (contract.SetOfContracts, error)
}

// contractSourcer provides validated source for requested contract/dir with contracts.
type contractSourcer interface {
	SourceInputs() []string
	SourceLoadType() types.SourceLoadType
	SourceFileType() types.SourceFileType
}

// Load allows to load contracts from provided source as `SetOfContracts`.
func Load(fs afero.Fs, source contractSourcer) (contracts contract.SetOfContracts, err error) {
	var (
		specificLoader = getLoader(fs, source)
		sourceType     = source.SourceLoadType()
	)

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

func getLoader(fs afero.Fs, s contractSourcer) contractLoader {
	switch s.SourceFileType() {
	case types.ProtoType:
		return proto.NewLoader(fs, s)

	case types.OpenAPIType:
		return openapi.NewLoader(fs, s)
	}

	return nil
}
