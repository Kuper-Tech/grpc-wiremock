package sourcer

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

var ErrContractTypeIsNotCorrespondToSourceLoadType = fmt.Errorf("expected another contract type")

// ErrProvidedDirectoryHasNoContracts indicates that provided directory input is invalid.
var ErrProvidedDirectoryHasNoContracts = fmt.Errorf("provided directory has no valid contracts")

// Sourcer hides how to work with contract sources.
//
// Contract source implements all details about opening, reading, validating
// files with specification. Such as .proto, .yaml (openapi) files.
type sourcer struct {
	contractPaths []string

	afero.Fs

	loadType types.SourceLoadType
	fileType types.SourceFileType
}

// New creates instance of sourcer.
func New(fs afero.Fs, path string, sourceFileType types.SourceFileType) (sourcer, error) {
	sourceLoadType, err := getContractSourceLoadType(fs, path)
	if err != nil {
		return sourcer{}, fmt.Errorf("get contract source load type: %w", err)
	}

	contractPaths, err := getContractPaths(fs, path, sourceFileType, sourceLoadType)
	if err != nil {
		return sourcer{}, fmt.Errorf("get contract paths: %w", err)
	}

	source := sourcer{
		Fs: fs,

		fileType: sourceFileType,
		loadType: sourceLoadType,

		contractPaths: contractPaths,
	}

	return source, nil
}

// SourceLoadType returns selected source load type.
func (c sourcer) SourceLoadType() types.SourceLoadType {
	return c.loadType
}

// SourceFileType returns selected source file type.
func (c sourcer) SourceFileType() types.SourceFileType {
	return c.fileType
}

// SourceInputs returns path of input entry.
//
// Input to the directory with the contracts or to the single contract file.
func (c sourcer) SourceInputs() []string {
	return c.contractPaths
}

func getContractSourceLoadType(fs afero.Fs, path string) (types.SourceLoadType, error) {
	entryInfo, err := fs.Stat(path)
	if err != nil {
		return types.UnknownLoadType, fmt.Errorf("stat entry %s: %w", path, err)
	}

	if entryInfo.IsDir() {
		return types.SourceDirType, nil
	}

	return types.SourceSingleType, nil
}

func getContractPaths(fs afero.Fs, path string, sourceFileType types.SourceFileType, loadType types.SourceLoadType) ([]string, error) {
	switch loadType {
	case types.SourceSingleType:
		if !sourceFileType.Is(fsutils.GetFileExt(path)) {
			return nil, ErrContractTypeIsNotCorrespondToSourceLoadType
		}

		return []string{path}, nil

	case types.SourceDirType:
		dirsWithContracts, err := findDirsWithContracts(fs, path, sourceFileType)
		if err != nil {
			return nil, fmt.Errorf("find dirs with contracts: %w", err)
		}

		if len(dirsWithContracts) == 0 {
			return nil, ErrProvidedDirectoryHasNoContracts
		}

		return dirsWithContracts, nil
	}

	return nil, fmt.Errorf("incorrect load type")
}

func findDirsWithContracts(fs afero.Fs, path string, sourceFileType types.SourceFileType) ([]string, error) {
	dirsWithContracts, err := fsutils.FindDirsWithContracts(fs, path, func(s string) bool {
		return sourceFileType.Is(fsutils.GetFileExt(s))
	})
	if err != nil {
		return nil, fmt.Errorf("find dirs: %w", err)
	}

	return dirsWithContracts, nil
}
