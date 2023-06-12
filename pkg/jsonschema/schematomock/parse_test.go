package schematomock

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/loader"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/mock"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

var (
	osfs = afero.NewOsFs()

	projectPath = filepath.Join(fsutils.CurrentDir(), "../../..")
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		fileType types.SourceFileType
		fs       afero.Fs
		want     []mock.Mock
	}{
		{
			fs: osfs, fileType: types.OpenAPIType,
			input: filepath.Join(projectPath, "static/tests/data/openapi/petstore-merged"),
			want:  []mock.Mock{{Description: "pet response", RequestUrlPath: "/pets", RequestMethod: "GET", ResponseStatusCode: 200}},
		},
		{
			fs: osfs, fileType: types.OpenAPIType,
			input: filepath.Join(projectPath, "static/tests/data/openapi/users-merged"),
			want:  []mock.Mock{{Description: "user response", RequestUrlPath: "/users", RequestMethod: "GET", ResponseStatusCode: 200}},
		},
		{
			fs: osfs, fileType: types.OpenAPIType,
			input: filepath.Join(projectPath, "static/tests/data/openapi/petstore-and-users-merged"),
			want: []mock.Mock{
				{Description: "user response", RequestUrlPath: "/users", RequestMethod: "GET", ResponseStatusCode: 200},
				{Description: "pet response", RequestUrlPath: "/pets", RequestMethod: "GET", ResponseStatusCode: 200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcerInstance, err := sourcer.New(tt.fs, tt.input, tt.fileType)
			require.NoError(t, err)

			contracts, err := loader.Load(tt.fs, sourcerInstance)
			require.NoError(t, err)

			descriptors, err := getDescriptors(contracts)
			require.NoError(t, err)

			var got []mock.Mock

			for _, desc := range descriptors {
				mocks, errParse := Parse(desc)
				require.NoError(t, errParse)

				got = append(got, mocks...)
			}

			require.Equal(t, len(got), len(tt.want))

			want := tt.want
			sort.Slice(want, func(i, j int) bool { return want[i].RequestUrlPath < want[j].RequestUrlPath })
			sort.Slice(got, func(i, j int) bool { return got[i].RequestUrlPath < got[j].RequestUrlPath })

			for idx := range tt.want {
				require.Equal(t, got[idx].Name, tt.want[idx].Name)
				require.Equal(t, got[idx].Description, tt.want[idx].Description)
				require.Equal(t, got[idx].RequestMethod, tt.want[idx].RequestMethod)
				require.Equal(t, got[idx].RequestUrlPath, tt.want[idx].RequestUrlPath)
				require.Equal(t, got[idx].ResponseStatusCode, tt.want[idx].ResponseStatusCode)
			}
		})
	}
}

func getDescriptors(contracts contract.SetOfContracts) ([]*openapi3.T, error) {
	descriptors := make([]*openapi3.T, 0, len(contracts))

	for _, contractIter := range contracts {
		desc, err := contract.OpenAPIFromAny(contractIter)
		if err != nil {
			return nil, fmt.Errorf("get desc: %w", err)
		}

		descriptors = append(descriptors, desc)
	}

	return descriptors, nil
}
