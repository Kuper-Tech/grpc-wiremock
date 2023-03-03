package fstesting

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting/entrymock"
)

func CreateMockFS(entries ...entrymock.Entry) afero.Fs {
	fs := afero.NewMemMapFs()

	// the MemMapFS write error is deliberately ignored
	fs, _ = WriteMockEntries(fs, entries...)

	return fs
}

func WriteMockEntries(fs afero.Fs, entries ...entrymock.Entry) (afero.Fs, error) {
	for _, entry := range entries {
		path := entry.Path()

		if entry.IsDir() {
			if err := fs.MkdirAll(path, os.ModePerm); err != nil {
				return nil, fmt.Errorf("create dir: %w", err)
			}
		} else {
			if err := fs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
				return nil, fmt.Errorf("create dir: %w", err)
			}
			_, err := fs.Create(path)
			if err != nil {
				return nil, fmt.Errorf("create file: %w", err)
			}
		}
	}

	return fs, nil
}
