package sourcer

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/fstesting/entrymock"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

func TestSource(t *testing.T) {
	testCases := []struct {
		name  string
		input string

		fileType types.SourceFileType

		fs afero.Fs

		expectedLoadType types.SourceLoadType
		expectedFileType types.SourceFileType
		expectedInputs   []string
	}{
		{
			input:          "/test/inner/dir_with_contracts",
			fileType:       types.ProtoType,
			expectedInputs: []string{"/test/inner/dir_with_contracts"},

			expectedLoadType: types.SourceDirType,
			expectedFileType: types.ProtoType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/test/inner/dir_with_contracts/sample.proto"),
			),
		},
		{
			input:          "/test/contract.proto",
			fileType:       types.ProtoType,
			expectedInputs: []string{"/test/contract.proto"},

			expectedLoadType: types.SourceSingleType,
			expectedFileType: types.ProtoType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/test/contract.proto"),
				entrymock.File("/test/inner/contract.proto"),
			),
		},
		{
			input:          "/test/openapi.yaml",
			fileType:       types.OpenAPIType,
			expectedInputs: []string{"/test/openapi.yaml"},

			expectedLoadType: types.SourceSingleType,
			expectedFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/test/openapi.yaml"),
			),
		},
		{
			input:            "/deps",
			fileType:         types.ProtoType,
			expectedInputs:   []string{"/deps/services/ph/grpc"},
			expectedLoadType: types.SourceDirType,
			expectedFileType: types.ProtoType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/deps/services/ab/openapi/ab.yaml"),
				entrymock.File("/deps/services/ph/grpc/phfd.proto"),
			),
		},
		{
			input:          "/deps",
			fileType:       types.OpenAPIType,
			expectedInputs: []string{"/deps/services"},

			expectedLoadType: types.SourceDirType,
			expectedFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/deps/services/ab.yaml"),
				entrymock.File("/deps/services/phfd.proto"),
			),
		},
		{
			input:          "/deps",
			fileType:       types.OpenAPIType,
			expectedInputs: []string{"/deps/services"},

			expectedLoadType: types.SourceDirType,
			expectedFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/deps/services/ab.yaml"),
				entrymock.File("/deps/services/phfd.proto"),
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			source, err := New(tt.fs, tt.input, tt.fileType)
			assert.NoError(t, err)

			assert.ElementsMatch(t, tt.expectedInputs, source.SourceInputs())
			assert.Equal(t, source.SourceLoadType(), tt.expectedLoadType)
			assert.Equal(t, source.SourceFileType(), tt.expectedFileType)
		})
	}
}

func TestSourceWithError(t *testing.T) {
	testCases := []struct {
		name  string
		input string

		fileType types.SourceFileType

		fs      afero.Fs
		wantErr error
	}{
		{
			name:    "dir doesn't exist",
			input:   "/test/inner/dir_with_contracts",
			wantErr: os.ErrNotExist,
			fs: fstesting.CreateMockFS(
				entrymock.Dir("/test/inner/dir"),
			),
		},
		{
			name:    "proto file doesn't exist",
			input:   "/test/contract_test.proto",
			wantErr: os.ErrNotExist,
			fs: fstesting.CreateMockFS(
				entrymock.File("/test/contract.proto"),
				entrymock.File("/test/inner/contract.proto"),
			),
		},
		{
			name:    "incorrect extension",
			input:   "/test/openapi.txt",
			wantErr: ErrContractTypeIsNotCorrespondToSourceLoadType,
			fs: fstesting.CreateMockFS(
				entrymock.File("/test/openapi.txt"),
			),
		},
		{
			name:    "empty dir",
			input:   "/test",
			wantErr: ErrProvidedDirectoryHasNoContracts,
			fs: fstesting.CreateMockFS(
				entrymock.Dir("/test"),
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.fs, tt.input, tt.fileType)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_findDirsWithContracts(t *testing.T) {
	tests := []struct {
		name           string
		fs             afero.Fs
		path           string
		sourceFileType types.SourceFileType
		want           []string
	}{
		{
			path:           "/deps",
			sourceFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.Dir("/deps/services/awesome/openapi/contract.yaml"),
				entrymock.Dir("/deps/services/awesome/grpc/contract.proto"),
				entrymock.Dir("/deps/services/disgusting/openapi/contract.yaml"),
			),
			want: []string{
				"/deps/services/awesome/openapi",
				"/deps/services/disgusting/openapi",
			},
		},
		{
			path:           "/deps",
			sourceFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.Dir("/deps/services/awesome/grpc/contract.yaml"),
				entrymock.Dir("/deps/services/awesome/grpc/contract.proto"),
				entrymock.Dir("/deps/services/disgusting/grpc/contract.yaml"),
			),
			want: []string{
				"/deps/services/awesome/grpc",
				"/deps/services/disgusting/grpc",
			},
		},
		{
			path:           "/deps",
			sourceFileType: types.OpenAPIType,
			fs: fstesting.CreateMockFS(
				entrymock.Dir("/deps/services/awesome/grpc/contract.proto"),
				entrymock.Dir("/deps/services/awesome/grpc/contract.proto"),
				entrymock.Dir("/deps/services/disgusting/grpc/contract.proto"),
			),
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findDirsWithContracts(tt.fs, tt.path, tt.sourceFileType)
			require.NoError(t, err)
			require.ElementsMatch(t, got, tt.want)
		})
	}
}
