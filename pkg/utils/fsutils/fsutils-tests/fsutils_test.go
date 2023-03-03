package fsutils_tests

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting/entrymock"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

func TestFindDirsWithContracts(t *testing.T) {
	tests := []struct {
		name string

		fs   afero.Fs
		path string

		want []string
	}{
		{
			path: "/searcher",
			fs: fstesting.CreateMockFS(
				entrymock.File("/searcher/services/navigation/grpc/navigation.proto"),
				entrymock.File("/searcher/services/product-hub/grpc/ph.proto"),
			),
			want: []string{
				"/searcher/services/navigation/grpc",
				"/searcher/services/product-hub/grpc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fsutils.FindDirsWithContracts(tt.fs, tt.path, func(entryName string) bool {
				return strings.Contains(entryName, ".proto") || strings.Contains(entryName, ".yaml")
			})
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestGatherEntries(t *testing.T) {
	tests := []struct {
		name string

		fs   afero.Fs
		path string

		want    []string
		wantErr error
	}{
		{
			path: "/",
			fs: fstesting.CreateMockFS(
				entrymock.File("/searcher/services/navigation/grpc/navigation.proto"),
				entrymock.File("/searcher/services/product-hub/grpc/ph.proto"),
			),
			want: []string{
				"/",
				"/searcher",
				"/searcher/services",
				"/searcher/services/navigation",
				"/searcher/services/navigation/grpc",
				"/searcher/services/navigation/grpc/navigation.proto",
				"/searcher/services/product-hub",
				"/searcher/services/product-hub/grpc",
				"/searcher/services/product-hub/grpc/ph.proto",
			},
		},
		{
			path: "/searcher/services/navigation/grpc",
			fs: fstesting.CreateMockFS(
				entrymock.File("/searcher/services/navigation/grpc/navigation.proto"),
				entrymock.File("/searcher/services/product-hub/grpc/ph.proto"),
			),
			want: []string{
				"/searcher/services/navigation/grpc",
				"/searcher/services/navigation/grpc/navigation.proto",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fsutils.GatherEntries(tt.fs, tt.path)
			assert.NoError(t, err)

			var paths []string
			for _, entry := range got {
				paths = append(paths, entry.Path)
			}
			assert.Equal(t, tt.want, paths)
		})
	}
}

func TestGetFileExt(t *testing.T) {
	tests := []struct {
		name string

		path string

		want string
	}{
		{path: "", want: ""},
		{path: "/foo/bar/file.txt", want: "txt"},
		{path: "/foo/bar/file.proto", want: "proto"},
		{path: "/foo/bar/file", want: ""},
		{path: "file.go", want: "go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fsutils.GetFileExt(tt.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemoveTmpDirs(t *testing.T) {
	tests := []struct {
		name string

		fs    afero.Fs
		paths []string

		wantFs afero.Fs
	}{
		{
			paths: []string{"/tmp/protos", "/tmp/foo"},
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
				entrymock.File("/searcher/services/navigation/grpc/navigation.proto"),
				entrymock.File("/searcher/services/product-hub/grpc/ph.proto"),
			),
			wantFs: fstesting.CreateMockFS(
				entrymock.Dir("/tmp/protos"),
				entrymock.Dir("/tmp/foo"),
				entrymock.File("/searcher/services/navigation/grpc/navigation.proto"),
				entrymock.File("/searcher/services/product-hub/grpc/ph.proto"),
			),
		},
		{
			paths: []string{"/tmp/path_does_not_exist"},
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
			wantFs: fstesting.CreateMockFS(
				entrymock.Dir("/tmp/path_does_not_exist"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsutils.RemoveTmpDirs(tt.fs, tt.paths...)
			assert.NoError(t, err)

			diff, err := fstesting.CompareFS(tt.wantFs, tt.fs)
			assert.NoError(t, err)
			assert.True(t, diff.Empty())
		})
	}
}

func TestRemoveWithSubdirs(t *testing.T) {
	tests := []struct {
		name string

		fs           afero.Fs
		path         string
		nameToRemove string

		wantFs afero.Fs
	}{
		{
			nameToRemove: ".keep",
			path:         "/",
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/.keep"),
				entrymock.File("/tmp/foo/bar/.keep"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
			wantFs: fstesting.CreateMockFS(
				entrymock.Dir("/tmp"),
				entrymock.Dir("/tmp/foo"),
				entrymock.Dir("/tmp/foo/bar"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
		},
		{
			nameToRemove: ".keep",
			path:         "/tmp/foo/bar",
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/.keep"),
				entrymock.File("/tmp/foo/bar/.keep"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
			wantFs: fstesting.CreateMockFS(
				entrymock.File("/tmp/.keep"),
				entrymock.Dir("/tmp/foo"),
				entrymock.Dir("/tmp/foo/bar"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsutils.RemoveWithSubdirs(tt.fs, tt.path, tt.nameToRemove)
			assert.NoError(t, err)

			diff, err := fstesting.CompareFS(tt.wantFs, tt.fs)
			assert.NoError(t, err)
			assert.True(t, diff.Empty())
		})
	}
}

func TestValidDirectories(t *testing.T) {
	tests := []struct {
		name string

		fs    afero.Fs
		paths []string

		wantErr error
	}{
		{
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/.keep"),
				entrymock.Dir("/tmp/foo"),
				entrymock.Dir("/tmp/foo/bar"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
			paths:   []string{"/tmp/foo", "/tmp/foo/bar"},
			wantErr: nil,
		},
		{
			fs: fstesting.CreateMockFS(
				entrymock.File("/tmp/.keep"),
				entrymock.Dir("/tmp/foo"),
				entrymock.Dir("/tmp/foo/bar"),
				entrymock.File("/tmp/protos/some_file.txt"),
				entrymock.File("/tmp/foo/some_file.proto"),
			),
			paths:   []string{"/tmp/foo/.keep", "/tmp/foo/bar"},
			wantErr: os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsutils.ValidDirectories(tt.fs, tt.paths...)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
