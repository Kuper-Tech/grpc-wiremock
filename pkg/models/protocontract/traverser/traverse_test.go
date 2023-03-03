package traverser

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract/loader"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

var projectDir = filepath.Join(fsutils.CurrentDir(), "../../../..")

func TestSetOfContracts_Descriptors(t *testing.T) {
	tests := []struct {
		name string

		fs   afero.Fs
		path string

		wantFiles []string
	}{
		{
			fs:   afero.NewOsFs(),
			path: filepath.Join(projectDir, "static/tests/data/static-examples/with-common-imports"),
			wantFiles: []string{
				"example.proto",
				"google/protobuf/empty.proto",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcerInstance, err := sourcer.New(tt.fs, tt.path, types.ProtoType)
			assert.NoError(t, err)

			contracts, err := loader.Load(tt.fs, sourcerInstance)
			assert.NoError(t, err)

			descriptors := Descriptors(contracts)
			assert.NotNil(t, descriptors)

			var files []string
			for _, desc := range descriptors {
				files = append(files, desc.GetName())
			}

			assert.ElementsMatch(t, tt.wantFiles, strutils.UniqueAndSorted(files...))
		})
	}
}
