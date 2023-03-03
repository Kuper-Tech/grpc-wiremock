package loader

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract/parser"
)

type loader struct {
	fs afero.Fs

	source contractSourcer
}

func NewLoader(fs afero.Fs, source contractSourcer) loader {
	return loader{fs: fs, source: source}
}

func (l loader) LoadFile() (protocontract.SetOfContracts, error) {
	protoFiles := l.source.SourceInputs()

	contracts, err := parser.ParseProtoFiles(protoFiles)
	if err != nil {
		return nil, fmt.Errorf("parse proto file: %w", err)
	}

	return contracts, nil
}

func (l loader) LoadDir() (protocontract.SetOfContracts, error) {
	protoPaths := l.source.SourceInputs()

	contracts, err := parser.ParseProtoDir(l.fs, protoPaths)
	if err != nil {
		return nil, fmt.Errorf("parse proto dir: %w", err)
	}

	return contracts, nil
}
