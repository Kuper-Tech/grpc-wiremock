package proto

import (
	"fmt"

	"github.com/spf13/afero"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/parser/proto"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
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
	protoFiles := l.source.SourceInputs()

	contracts, err := proto.ParseProtoFiles(protoFiles)
	if err != nil {
		return nil, fmt.Errorf("parse proto file: %w", err)
	}

	return contracts, nil
}

func (l loader) LoadDir() (contract.SetOfContracts, error) {
	protoPaths := l.source.SourceInputs()

	contracts, err := proto.ParseProtoDir(l.fs, protoPaths)
	if err != nil {
		return nil, fmt.Errorf("parse proto dir: %w", err)
	}

	return contracts, nil
}
