package entrymock

import (
	"io/fs"
	"path/filepath"
	"time"
)

type Entry struct {
	isDir bool
	name  string
	path  string
}

func Dir(path string) Entry {
	return NewEntryMock(path, true)
}

func File(path string) Entry {
	return NewEntryMock(path, false)
}

func NewEntryMock(path string, isDir bool) Entry {
	return Entry{name: filepath.Base(path), isDir: isDir, path: path}
}

func (e Entry) Name() string {
	return e.name
}

func (e Entry) IsDir() bool {
	return e.isDir
}

func (e Entry) Path() string {
	return e.path
}

func (e Entry) Size() int64 {
	return 0
}

func (e Entry) Mode() fs.FileMode {
	return fs.ModePerm
}

func (e Entry) ModTime() time.Time {
	return time.Time{}
}

func (e Entry) Sys() any {
	return nil
}
