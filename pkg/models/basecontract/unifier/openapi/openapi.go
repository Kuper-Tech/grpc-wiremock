package openapi

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type openapiUnifier struct {
	fs afero.Fs
}

func NewUnifier(fs afero.Fs) *openapiUnifier {
	return &openapiUnifier{fs: fs}
}

func (o *openapiUnifier) Unify(ctx context.Context, contracts contract.SetOfContracts, path string) error {
	for _, contractToUnify := range contracts {
		if err := o.saveFromDescriptor(contractToUnify, path); err != nil {
			return fmt.Errorf("save from descriptor: %w", err)
		}
	}

	return nil
}

func (o *openapiUnifier) saveFromDescriptor(contractToUnify contract.Contract, path string) error {
	descriptor, err := contract.OpenAPIFromAny(contractToUnify)
	if err != nil {
		return fmt.Errorf("get openapi descriptor: %w", err)
	}

	content, err := yaml.Marshal(descriptor)
	if err != nil {
		return fmt.Errorf("marshal contract %s: %w", contractToUnify.HeaderPath, err)
	}

	pathToSave := filepath.Join(path, filepath.Base(contractToUnify.HeaderPath))

	if err = fsutils.WriteFile(o.fs, pathToSave, string(content)); err != nil {
		return fmt.Errorf("write unified contract %s: %w", contractToUnify.HeaderPath, err)
	}

	return nil
}
