package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/proxy"
)

func Test_replaceValue(t *testing.T) {
	tests := []struct {
		name string

		obj proxy.SubstitutionForMainGoView

		want interface{}
	}{
		{
			obj: proxy.SubstitutionForMainGoView{
				GoPackages: []string{"github.com/service"},

				OriginalGoPackagesWithService: []string{"service1", "service2"},
				PackageToServices: []proxy.PackageToService{
					{ProtoPackage: "service", Service: "service1"},
					{ProtoPackage: "service", Service: "service2"},
				},
			},
			want: proxy.SubstitutionForMainGoView{
				GoPackages: []string{"github.com/service"},

				OriginalGoPackagesWithService: []string{"service1", "service2"},
				PackageToServices: []proxy.PackageToService{
					{ProtoPackage: "service", Service: "service1"},
					{ProtoPackage: "service", Service: "service2"},
				},
			},
		},
		{
			obj: proxy.SubstitutionForMainGoView{
				GoPackages: []string{"github.com/service"},

				OriginalGoPackagesWithService: []string{"go_resolved", "service2"},
				PackageToServices: []proxy.PackageToService{
					{ProtoPackage: "type", Service: "service1"},
					{ProtoPackage: "service", Service: "service2"},
				},
			},
			want: proxy.SubstitutionForMainGoView{
				GoPackages: []string{"github.com/service"},

				OriginalGoPackagesWithService: []string{"go_resolved", "service2"},
				PackageToServices: []proxy.PackageToService{
					{ProtoPackage: "type_resolved", Service: "service1"},
					{ProtoPackage: "service", Service: "service2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolveCollisions(&tt.obj)
			assert.Equal(t, tt.want, tt.obj)
		})
	}
}
