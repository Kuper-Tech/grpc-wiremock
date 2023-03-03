package openapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/afero"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

func ParseOpenAPIFiles(openapiFile []string) (contract.SetOfContracts, error) {
	var contracts contract.SetOfContracts

	contractPath := sliceutils.FirstOf(openapiFile)

	descriptor, err := parse(contractPath)
	if err != nil {
		return nil, fmt.Errorf("parse openapi: %w", err)
	}

	convertedContract, errConvert := contract.FromOpenAPIDescriptor(descriptor, contractPath)
	if errConvert != nil {
		return nil, fmt.Errorf("convert openapi: %w", errConvert)
	}

	contracts = append(contracts, convertedContract)

	return contracts, nil
}

func ParseOpenAPIDir(fs afero.Fs, openapiPath string) (contract.SetOfContracts, error) {
	var contracts contract.SetOfContracts

	files, err := fsutils.GatherMatchedEntriesInDir(fs, openapiPath, onlyOpenAPIFiles)
	if err != nil {
		return nil, fmt.Errorf("get valid files: %w", err)
	}

	for _, openapiFile := range files {
		descriptor, err := parse(openapiFile)
		if err != nil {
			return nil, fmt.Errorf("parse openapi: %w", err)
		}

		convertedContract, errConvert := contract.FromOpenAPIDescriptor(descriptor, openapiFile)
		if errConvert != nil {
			return nil, fmt.Errorf("convert openapi: %w", errConvert)
		}

		contracts = append(contracts, convertedContract)
	}

	return contracts, nil
}

func parse(openapiFile string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	descriptor, err := loader.LoadFromFile(openapiFile)
	if err != nil {
		return nil, fmt.Errorf("load from data '%s': %w", openapiFile, err)
	}

	ctx := context.Background()
	contractPath := filepath.Dir(openapiFile)
	internalizeF := createInternalizeF(contractPath)

	descriptor.InternalizeRefs(ctx, internalizeF)

	return descriptor, nil
}

func createInternalizeF(contractPath string) func(string) string {
	return func(ref string) string {
		if ref == "" {
			return ""
		}

		split := strings.SplitN(ref, "#", 2)
		if len(split) == 2 {
			return filepath.Base(split[1])
		}

		ref = sliceutils.FirstOf(split)

		for ext := filepath.Ext(ref); ext != ""; ext = filepath.Ext(ref) {
			ref = strings.TrimSuffix(ref, ext)
		}

		ref = strings.TrimPrefix(strings.TrimPrefix(ref, contractPath), "/")

		return strutils.ToCamelCase(ref)
	}
}

func onlyOpenAPIFiles(info os.FileInfo) bool {
	if !info.IsDir() &&
		strings.Contains(info.Name(), ".yaml") {
		return true
	}
	return false
}
