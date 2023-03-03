package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

func Test_createDataToSubstituteMainGo(t *testing.T) {
	tests := []struct {
		name string

		contracts contract.SetOfContracts

		wantData SubstitutionForMainGoView
	}{
		{
			contracts: contract.SetOfContracts{
				contract.Contract{
					Base: contract.Base{
						Package:   "example",
						GoPackage: "gitlab.com/paas/example",
						Services: []contract.Service{
							{
								Name: "Example",
								Methods: []contract.Method{
									{
										Name:       "Unary",
										Package:    "example",
										MethodType: types.UnaryType,
									},
								},
							},
						},
					},
				},
				contract.Contract{
					Base: contract.Base{
						Package:   "example",
						GoPackage: "gitlab.com/paas/example",
						Services: []contract.Service{
							{
								Name: "AnotherExample",
								Methods: []contract.Method{
									{
										Name:       "AnotherUnary",
										Package:    "example",
										MethodType: types.UnaryType,
									},
								},
							},
						},
					},
				},
			},
			wantData: SubstitutionForMainGoView{
				OriginalGoPackagesWithService: []string{
					"gitlab.com/paas/example/another_example",
					"gitlab.com/paas/example/example",
				},
				GoPackages: []string{
					"gitlab.com/paas/example",
				},
				PackageToServices: []PackageToService{
					{ProtoPackage: "example", Service: "Example", ServicePackage: "example_example"},
					{ProtoPackage: "example", Service: "AnotherExample", ServicePackage: "example_another_example"},
				},
			},
		},
		{
			contracts: contract.SetOfContracts{
				contract.Contract{
					Base: contract.Base{
						Package:   "example",
						GoPackage: "gitlab.com/paas/example",
						Services: []contract.Service{
							{
								Name: "Example",
							},
						},
					},
				},
				contract.Contract{
					Base: contract.Base{
						Package:   "example",
						GoPackage: "gitlab.com/paas/example",
						Services: []contract.Service{
							{
								Name: "AnotherExample",
							},
						},
					},
				},
			},
			wantData: SubstitutionForMainGoView{
				OriginalGoPackagesWithService: nil,
				GoPackages:                    nil,
				PackageToServices:             nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := substitutionForProject(tt.contracts)
			assert.Equal(t, tt.wantData, got)
		})
	}
}

func Test_substitutionMethodForStubs(t *testing.T) {
	tests := []struct {
		name string

		method  contract.Method
		service contract.Service
		baseURL string

		want SubstitutionMethodForStubsView
	}{
		{
			service: contract.Service{
				Name:      "Example",
				GoPackage: "gitlab.io/paas/example",
			},
			baseURL: "http://localhost:8000",
			method: contract.Method{
				Name:       "Unary",
				Package:    "example",
				InType:     contract.MessageType{Name: "CustomEmpty", Package: "import", GoPackage: "gitlab.io/paas/example/import", IsExternal: true},
				OutType:    contract.MessageType{Name: "Response", Package: "example", GoPackage: "gitlab.io/paas/example", IsExternal: false},
				MethodType: types.UnaryType,
			},
			want: SubstitutionMethodForStubsView{
				Method:           "Unary",
				MethodInPackage:  "import",
				Package:          "example",
				Service:          "Example",
				MethodPackage:    "example",
				MethodOutPackage: "example",
				MethodOutName:    "Response",
				MethodInName:     "CustomEmpty",
				PackageHeader:    "example_example",
				URL:              "http://localhost:8000/Example/Unary",
				GoPackages:       []string{"gitlab.io/paas/example", "gitlab.io/paas/example/import"},
			},
		},
		{
			service: contract.Service{
				Name:      "Fancy",
				GoPackage: "gitlab.io/paas/example",
			},
			baseURL: "http://localhost:8001",
			method: contract.Method{
				Name:       "SomeFancyName",
				Package:    "fancy",
				InType:     contract.MessageType{Name: "CustomEmpty", Package: "another", GoPackage: "gitlab.io/paas/another", IsExternal: true},
				OutType:    contract.MessageType{Name: "ResponseFromExternalPackage", Package: "external", GoPackage: "gitlab.io/paas/external", IsExternal: true},
				MethodType: types.UnaryType,
			},
			want: SubstitutionMethodForStubsView{
				Service:          "Fancy",
				MethodPackage:    "fancy",
				Package:          "example",
				MethodInPackage:  "another",
				MethodOutPackage: "external",
				MethodInName:     "CustomEmpty",
				Method:           "SomeFancyName",
				PackageHeader:    "example_fancy",
				MethodOutName:    "ResponseFromExternalPackage",
				URL:              "http://localhost:8001/Fancy/SomeFancyName",
				GoPackages:       []string{"gitlab.io/paas/another", "gitlab.io/paas/external"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := substitutionMethodForStubs(tt.service, tt.method, tt.baseURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
