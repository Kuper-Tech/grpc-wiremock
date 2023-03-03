package fsutils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

func MakeZipArchive(sourceFS, targetFS afero.Fs, sourcePath, targetPath string) error {
	f, err := targetFS.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	waklF := createWalkF(sourceFS, sourcePath, writer)

	if err := afero.Walk(sourceFS, sourcePath, waklF); err != nil {
		return fmt.Errorf("walk: %w", err)
	}

	return nil
}

func createWalkF(sourceFS afero.Fs, sourcePath string, writer *zip.Writer) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk func: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		header, err := createHeader(info, sourcePath, path)
		if err != nil {
			return fmt.Errorf("create header: %w", err)
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		f, err := sourceFS.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}

		return nil
	}
}

func createHeader(info fs.FileInfo, sourcePath, path string) (*zip.FileHeader, error) {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, fmt.Errorf("header: %w", err)
	}

	header.Method = zip.Deflate

	header.Name, err = filepath.Rel(filepath.Dir(sourcePath), path)
	if err != nil {
		return nil, fmt.Errorf("get rel path: %w", err)
	}

	if info.IsDir() {
		header.Name += "/"
	}

	return header, nil
}
