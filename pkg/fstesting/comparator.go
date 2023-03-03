package fstesting

import (
	"bytes"
	"fmt"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting/entry"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type FSDiff struct {
	EntriesWithNotEqualBodies []entry.Entry
}

func (diff *FSDiff) Add(entry entry.Entry) {
	diff.EntriesWithNotEqualBodies =
		append(diff.EntriesWithNotEqualBodies, entry)
}

func (diff *FSDiff) Empty() bool {
	return len(diff.EntriesWithNotEqualBodies) == 0
}

type FSEntries []entry.Entry

type SetOfFSEntries struct {
	Entries []FSEntries

	ActualFS   afero.Fs
	ExpectedFS afero.Fs
}

func (sfse *SetOfFSEntries) Expected() FSEntries {
	if len(sfse.Entries) < 2 {
		panic("set is incompatible")
	}
	return sfse.Entries[0]
}

func (sfse *SetOfFSEntries) Actual() FSEntries {
	if len(sfse.Entries) < 2 {
		panic("set is incompatible")
	}
	return sfse.Entries[1]
}

func (sfse *SetOfFSEntries) Size() int {
	return len(sfse.Actual())
}

func NewSet(expectedFS, actualFS afero.Fs) (SetOfFSEntries, error) {
	expectedEntries, err := fsutils.GatherEntries(expectedFS, "/")
	if err != nil {
		return SetOfFSEntries{}, fmt.Errorf("walk throw expected fs: %w", err)
	}

	actualEntries, err := fsutils.GatherEntries(actualFS, "/")
	if err != nil {
		return SetOfFSEntries{}, fmt.Errorf("walk throw actual fs: %w", err)
	}

	if len(expectedEntries) != len(actualEntries) {
		return SetOfFSEntries{}, fmt.Errorf(
			"incorrect size of fs entries: len(expected) = %d, len(actual) = %d",
			len(expectedEntries), len(actualEntries),
		)
	}

	return SetOfFSEntries{
		ActualFS:   actualFS,
		ExpectedFS: expectedFS,
		Entries:    []FSEntries{expectedEntries, actualEntries},
	}, nil
}

func (sfse *SetOfFSEntries) compareFiles(expected, actual entry.Entry) (bool, error) {
	entryIsDir := expected.FileInfo.IsDir() &&
		expected.FileInfo.IsDir() == actual.FileInfo.IsDir()

	if entryIsDir {
		return true, nil
	}

	expectedBody, err := fsutils.ReadFile(sfse.ExpectedFS, expected.Path)
	if err != nil {
		return false, fmt.Errorf("read expected file: %w", err)
	}

	actualBody, err := fsutils.ReadFile(sfse.ActualFS, actual.Path)
	if err != nil {
		return false, fmt.Errorf("read actual file: %w", err)
	}

	equal := bytes.Equal(expectedBody, actualBody)

	return equal, nil
}

func (sfse *SetOfFSEntries) Compare() (FSDiff, error) {
	var (
		size     = sfse.Size()
		actual   = sfse.Actual()
		expected = sfse.Expected()
	)

	var diff FSDiff
	for i := 0; i < size; i++ {
		var (
			actualEntry   = actual[i]
			expectedEntry = expected[i]
		)

		equal, err := sfse.compareFiles(expectedEntry, actualEntry)
		if err != nil {
			return FSDiff{}, fmt.Errorf("compare files: %w", err)
		}

		if !equal {
			diff.Add(actualEntry)
		}
	}

	return diff, nil
}

func CompareFS(expectedFS, actualFS afero.Fs) (FSDiff, error) {
	sfse, err := NewSet(expectedFS, actualFS)
	if err != nil {
		return FSDiff{}, fmt.Errorf("create set of fs: %w", err)
	}

	fsDiff, err := sfse.Compare()
	if err != nil {
		return FSDiff{}, fmt.Errorf("compare fs: %w", err)
	}

	return fsDiff, nil
}
