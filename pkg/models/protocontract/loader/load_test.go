package loader

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

var projectPath = filepath.Join(fsutils.CurrentDir(), "../../../..")

var (
	osfs = afero.NewOsFs()
)

func Test_loader_loadFromDir(t *testing.T) {
	tests := []struct {
		name string

		path         string
		fs           afero.Fs
		contractType types.SourceFileType

		want    protocontract.SetOfContracts
		wantErr error
	}{
		{
			fs: osfs, contractType: types.ProtoType,
			path: filepath.Join(projectPath, "static/tests/data/static-examples/with-local-imports-as-reciever"),
			want: protocontract.SetOfContracts{
				protocontract.Contract{
					Base: struct {
						Package      string
						GoPackage    string
						HeaderPath   string
						Messages     []string
						ImportsPaths []string
						Services     []protocontract.Service
					}{
						Package:      "example",
						GoPackage:    "gitlab.io/paas/example",
						HeaderPath:   filepath.Join(projectPath, "static/tests/data/static-examples/with-local-imports-as-reciever/api/grpc/example.proto"),
						ImportsPaths: []string{filepath.Join(projectPath, "static/tests/data/static-examples/with-local-imports-as-reciever/api/grpc/local/custom.proto")},
						Messages:     []string{"Request", "Response"},
						Services: []protocontract.Service{
							{
								Name:      "Example",
								GoPackage: "gitlab.io/paas/example",
								Methods: []protocontract.Method{
									{
										Name:       "Unary",
										Package:    "example",
										InType:     protocontract.MessageType{Name: "CustomEmpty", Package: "local", GoPackage: "gitlab.io/paas/example/local", IsExternal: true},
										OutType:    protocontract.MessageType{Name: "Response", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										MethodType: types.UnaryType,
									},
									{
										Name:       "ClientSideStream",
										Package:    "example",
										InType:     protocontract.MessageType{Name: "Request", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										OutType:    protocontract.MessageType{Name: "Response", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										MethodType: types.ClientSideStreamingType,
									},
									{
										Name:       "ServerSideStream",
										Package:    "example",
										InType:     protocontract.MessageType{Name: "Request", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										OutType:    protocontract.MessageType{Name: "Response", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										MethodType: types.ServerSideStreamingType,
									},
									{
										Name:       "BidirectionalStream",
										Package:    "example",
										InType:     protocontract.MessageType{Name: "Request", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										OutType:    protocontract.MessageType{Name: "Response", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
										MethodType: types.BidirectionalStreamingType,
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			fs: osfs, contractType: types.OpenAPIType,
			path:    filepath.Join(projectPath, "static/tests/data/openapi/petstore"),
			want:    nil,
			wantErr: unsupportedLoaderErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcerInstance, err := sourcer.New(tt.fs, tt.path, tt.contractType)
			assert.NoError(t, err)

			contracts, err := Load(tt.fs, sourcerInstance)
			assert.ErrorIs(t, err, tt.wantErr)

			var contractsToCompare protocontract.SetOfContracts
			for _, con := range contracts {
				contractsToCompare = append(contractsToCompare, protocontract.Contract{Base: con.Base})
			}

			assert.Equal(t, tt.want, contractsToCompare)
		})
	}
}
