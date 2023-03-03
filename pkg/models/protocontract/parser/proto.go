package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

func ParseProtoFiles(protoFiles []string) (protocontract.SetOfContracts, error) {
	var contracts protocontract.SetOfContracts

	protoPath := filepath.Dir(sliceutils.FirstOf(protoFiles))

	descriptors, err := parse([]string{protoPath, environment.TmpWellKnownProtosDir}, protoFiles...)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}

	for _, descriptor := range descriptors {
		convertedContract, errConvert := protocontract.FromProtoDescriptor(descriptor, protoPath)
		if errConvert != nil {
			return nil, fmt.Errorf("convert proto: %w", errConvert)
		}

		contracts = append(contracts, convertedContract)
	}

	return contracts, nil
}

func ParseProtoDir(fs afero.Fs, protoPaths []string) (protocontract.SetOfContracts, error) {
	var contracts protocontract.SetOfContracts

	for _, protoPath := range protoPaths {
		files, err := fsutils.GatherMatchedEntriesInDir(fs, protoPath, onlyProtoFiles)
		if err != nil {
			return nil, fmt.Errorf("get valid files: %w", err)
		}

		descriptors, err := parse([]string{protoPath, environment.TmpWellKnownProtosDir}, files...)
		if err != nil {
			return nil, fmt.Errorf("parse proto: %w", err)
		}

		for _, descriptor := range descriptors {
			convertedContract, errConvert := protocontract.FromProtoDescriptor(descriptor, protoPath)
			if errConvert != nil {
				return nil, fmt.Errorf("convert proto: %w", errConvert)
			}

			contracts = append(contracts, convertedContract)
		}
	}

	return contracts, nil
}

func parse(protoPaths []string, protoFiles ...string) ([]*desc.FileDescriptor, error) {
	resolvedHeaders, err := protoparse.ResolveFilenames(protoPaths, protoFiles...)
	if err != nil {
		return nil, fmt.Errorf("resolve names: %w", err)
	}

	parser := &protoparse.Parser{
		IncludeSourceCodeInfo: true,
		ImportPaths:           protoPaths,
	}

	descriptors, err := parser.ParseFiles(resolvedHeaders...)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	return descriptors, nil
}

func onlyProtoFiles(info os.FileInfo) bool {
	if !info.IsDir() &&
		strings.Contains(info.Name(), ".proto") {
		return true
	}
	return false
}
