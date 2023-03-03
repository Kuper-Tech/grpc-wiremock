package openapi

import (
	"fmt"

	"github.com/spf13/afero"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/parser/openapi"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

type contractSourcer interface {
	SourceInputs() []string
	SourceLoadType() types.SourceLoadType
	SourceFileType() types.SourceFileType
}

type loader struct {
	fs     afero.Fs
	source contractSourcer
}

func NewLoader(fs afero.Fs, source contractSourcer) loader {
	return loader{fs: fs, source: source}
}

func (l loader) LoadFile() (contract.SetOfContracts, error) {
	oaFiles := l.source.SourceInputs()

	contracts, err := openapi.ParseOpenAPIFiles(oaFiles)
	if err != nil {
		return nil, fmt.Errorf("parse openapi file: %w", err)
	}

	return contracts, nil
}

func (l loader) LoadDir() (contract.SetOfContracts, error) {
	oaFiles := l.source.SourceInputs()

	contracts, err := openapi.ParseOpenAPIDir(l.fs, sliceutils.FirstOf(oaFiles))
	if err != nil {
		return nil, fmt.Errorf("parse openapi dir: %w", err)
	}

	return contracts, nil
}
