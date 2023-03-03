package fsutils

import (
	"bytes"
	"fmt"
	stdfs "io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting/entry"
)

func CurrentDir() string {
	_, filePath, _, _ := runtime.Caller(1)
	return filepath.Dir(filePath)
}

func FindDirsWithContracts(fs afero.Fs, path string, match func(string) bool) ([]string, error) {
	var dirsWithContracts []string

	walkF := func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk: %w", err)
		}

		if info.IsDir() {
			entries, err := afero.ReadDir(fs, path)
			if err != nil {
				return fmt.Errorf("read dir: %w", err)
			}

			for _, ent := range entries {
				if match(ent.Name()) {
					dirsWithContracts = append(dirsWithContracts, path)
					return filepath.SkipDir
				}
			}
		}

		return nil
	}

	if err := afero.Walk(fs, path, walkF); err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}

	return dirsWithContracts, nil
}

func GatherEntries(fs afero.Fs, path string) ([]entry.Entry, error) {
	var entries []entry.Entry

	walkF := func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walkF, path '%s': %w", path, err)
		}

		entries = append(entries, entry.NewEntry(path, info))

		return nil
	}

	if err := afero.Walk(fs, path, walkF); err != nil {
		return nil, fmt.Errorf("walk throw fs: %w", err)
	}

	return entries, nil
}

func ReadFile(fs afero.Fs, path string) ([]byte, error) {
	body, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("read file %s body: %w", path, err)
	}

	return bytes.TrimSpace(body), nil
}

func ValidDirectories(fs afero.Fs, paths ...string) error {
	for _, path := range paths {
		info, err := fs.Stat(path)
		if err != nil {
			return fmt.Errorf("stat '%s': %w", path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("is not dir '%s': %w", path, err)
		}
	}
	return nil
}

func GetFileExt(path string) string {
	return strings.TrimPrefix(filepath.Ext(filepath.Base(path)), ".")
}

func WriteFile(fs afero.Fs, path string, content string) error {
	dir := filepath.Dir(path)

	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("create dir: %s, err: %w", dir, err)
	}

	if err := afero.WriteFile(fs, path, []byte(content), os.ModePerm); err != nil {
		return fmt.Errorf("write file: %s, err: %w", path, err)
	}

	return nil
}

func RemoveTmpDirs(fs afero.Fs, paths ...string) error {
	for _, path := range paths {
		if err := fs.RemoveAll(path); err != nil {
			return fmt.Errorf("remove tmp dir '%s', %w", path, err)
		}
		if err := fs.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("create tmp dir '%s', %w", path, err)
		}
	}

	return nil
}

func RemoveWithSubdirs(fs afero.Fs, path, name string) error {
	walkF := func(path string, info stdfs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walkF: %w", err)
		}

		if info.IsDir() || info.Name() != name {
			return nil
		}

		if err = fs.Remove(path); err != nil {
			return fmt.Errorf("remove: %w", err)
		}

		return nil
	}

	if err := afero.Walk(fs, path, walkF); err != nil {
		return fmt.Errorf("walk throw fs: %w", err)
	}

	return nil
}

func GatherDirs(fs afero.Fs, projectsDir string) ([]os.FileInfo, error) {
	entries, err := afero.ReadDir(fs, projectsDir)
	if err != nil {
		return nil, fmt.Errorf("read dir with projects: %w", err)
	}

	var filtered []os.FileInfo
	for _, ent := range entries {
		if !ent.IsDir() || IsHiddenFile(ent.Name()) {
			continue
		}

		filtered = append(filtered, ent)
	}

	return filtered, nil
}

func IsHiddenFile(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

func GatherMatchedEntriesInDir(fs afero.Fs, path string, match func(info os.FileInfo) bool) ([]string, error) {
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var headers []string
	for _, entry := range entries {
		if match(entry) {
			path := filepath.Join(path, entry.Name())
			headers = append(headers, path)
		}
	}

	return headers, nil
}
