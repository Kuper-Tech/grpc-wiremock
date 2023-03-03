package fsutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func CopyDir(sourceFS, targetFS afero.Fs, sourcePath, targetPath string, skipIfExists bool) error {
	info, err := sourceFS.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("stat dir: %w", err)
	}
	return copy(sourceFS, targetFS, sourcePath, targetPath, info, skipIfExists)
}

func copy(sourceFS, targetFS afero.Fs, src, dest string, info os.FileInfo, skipIfExists bool) error {
	if info.Mode()&os.ModeDevice != 0 {
		return nil
	}

	switch {
	case info.IsDir():
		if err := dcopy(sourceFS, targetFS, src, dest, skipIfExists); err != nil {
			return fmt.Errorf("copy dir: %w", err)
		}
	default:
		if err := fcopy(sourceFS, targetFS, src, dest); err != nil {
			return fmt.Errorf("copy file: %w", err)
		}
	}

	return nil
}

func fcopy(sourceFS, targetFS afero.Fs, src, dest string) error {
	if err := targetFS.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return fmt.Errorf("make dirs: %w", err)
	}

	f, err := targetFS.Create(dest)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	s, err := sourceFS.Open(src)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	var buf []byte = nil
	var (
		w io.Writer = f
		r io.Reader = s
	)

	if _, err = io.CopyBuffer(w, r, buf); err != nil {
		return fmt.Errorf("copy buffer: %w", err)
	}

	return nil
}

func dcopy(sourceFS, targetFS afero.Fs, sourcePath, targetPath string, skipIfExists bool) (err error) {
	if skip, err := onDirExists(targetFS, targetPath, skipIfExists); err != nil {
		return err
	} else if skip {
		return nil
	}

	contents, err := afero.ReadDir(sourceFS, sourcePath)
	if err != nil {
		return
	}

	for _, content := range contents {
		cs, cd := filepath.Join(sourcePath, content.Name()), filepath.Join(targetPath, content.Name())

		if err = copy(sourceFS, targetFS, cs, cd, content, skipIfExists); err != nil {
			// If any error, exit immediately
			return
		}
	}

	return
}

func onDirExists(targetFS afero.Fs, targetPath string, skipIfExists bool) (bool, error) {
	_, err := targetFS.Stat(targetPath)
	if err == nil {
		if skipIfExists {
			return false, nil
		}
		if err = targetFS.RemoveAll(targetPath); err != nil {
			return false, err
		}
	} else if err != nil && !os.IsNotExist(err) {
		return true, err // Unwelcome error type...!
	}
	return false, nil
}
