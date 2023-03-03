package proto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

var ProtoPaths = []string{
	environment.TmpWellKnownProtosDir,
	environment.TmpAnnotationProtosDir,
}

func ParseProtoFiles(protoFiles []string) (basecontract.SetOfContracts, error) {
	contractProtoPath := filepath.Dir(sliceutils.FirstOf(protoFiles))

	protoPaths := []string{contractProtoPath}
	protoPaths = append(protoPaths, ProtoPaths...)

	descriptors, err := parse(protoPaths, protoFiles...)
	if err != nil {
		return nil, fmt.Errorf("parse proto: %w", err)
	}

	var contracts basecontract.SetOfContracts

	for _, descriptor := range descriptors {
		convertedContract, errConvert := basecontract.FromProtoDescriptor(descriptor, contractProtoPath)
		if errConvert != nil {
			return nil, fmt.Errorf("convert proto: %w", errConvert)
		}

		contracts = append(contracts, convertedContract)
	}

	return contracts, nil
}

func ParseProtoDir(fs afero.Fs, protoPaths []string) (basecontract.SetOfContracts, error) {
	var contracts basecontract.SetOfContracts
	for _, protoPath := range protoPaths {
		files, err := fsutils.GatherMatchedEntriesInDir(fs, protoPath, onlyProtoFiles)
		if err != nil {
			return nil, fmt.Errorf("get valid files: %w", err)
		}

		protoPaths := []string{protoPath}
		protoPaths = append(protoPaths, ProtoPaths...)

		descriptors, err := parse(protoPaths, files...)
		if err != nil {
			return nil, fmt.Errorf("parse proto: %w", err)
		}

		for _, descriptor := range descriptors {
			convertedContract, errConvert := basecontract.FromProtoDescriptor(descriptor, protoPath)
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
