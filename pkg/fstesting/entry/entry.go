package entry

import "io/fs"

type Entry struct {
	Path string

	fs.FileInfo
}

func NewEntry(path string, info fs.FileInfo) Entry {
	return Entry{FileInfo: info, Path: path}
}
