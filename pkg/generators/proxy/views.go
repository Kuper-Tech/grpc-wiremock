package proxy

import (
	"fmt"
	"strings"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

type PackageToService struct {
	Service string

	ProtoPackage   string
	ServicePackage string
}

// SubstitutionForMainGoView data view for generating 'main.go'.
type SubstitutionForMainGoView struct {
	GoPackages                    []string
	OriginalGoPackagesWithService []string
	PackageToServices             []PackageToService
}

// SubstitutionMethodForStubsView data view for generating file with grpc method.
type SubstitutionMethodForStubsView struct {
	URL              string
	Method           string
	Service          string
	PackageHeader    string
	Package          string
	MethodInName     string
	MethodInPackage  string
	MethodOutName    string
	MethodOutPackage string
	MethodPackage    string

	GoPackages []string
}

// SubstitutionServiceForStubsView data view for generating file with service.
type SubstitutionServiceForStubsView struct {
	Service   string
	Package   string
	GoPackage string

	PackageHeader string
}

func substitutionForProject(contracts contract.SetOfContracts) SubstitutionForMainGoView {
	data := SubstitutionForMainGoView{}

	gatherPackageToServices := func(contractDesc contract.Contract) {
		for _, service := range contractDesc.Services {
			packageName := contractDesc.Package

			data.PackageToServices = append(
				data.PackageToServices,
				PackageToService{
					Service:        strings.Title(service.Name),
					ProtoPackage:   strutils.ToSnakeCase(packageName),
					ServicePackage: strutils.ToSnakeCase(strutils.ToPackageName(contractDesc.GoPackage) + service.Name),
				},
			)
		}
	}

	gatherOriginalGoPackagesWithService := func(contractDesc contract.Contract) {
		for _, service := range contractDesc.Services {
			data.OriginalGoPackagesWithService = append(
				data.OriginalGoPackagesWithService,
				fmt.Sprintf(
					"%s/%s",
					toOriginalPackageName(contractDesc.GoPackage),
					strutils.ToSnakeCase(service.Name),
				),
			)
		}
	}

	for _, contractDesc := range contracts.WithMethodsOnly() {
		gatherPackageToServices(contractDesc)
		gatherOriginalGoPackagesWithService(contractDesc)
		data.GoPackages = append(data.GoPackages, contractDesc.GoPackage)
	}

	data.GoPackages = strutils.UniqueAndSorted(data.GoPackages...)
	data.OriginalGoPackagesWithService = strutils.UniqueAndSorted(data.OriginalGoPackagesWithService...)

	return data
}

func substitutionMethodForStubs(service contract.Service, method contract.Method, baseURL string) (SubstitutionMethodForStubsView, error) {
	var goPackages []string

	switch method.MethodType {
	case types.UnaryType:
		goPackages = strutils.UniqueAndSorted(method.InType.GoPackage, method.OutType.GoPackage)
	case types.ServerSideStreamingType:
		goPackages = strutils.UniqueAndSorted(method.InType.GoPackage, method.OutType.GoPackage, service.GoPackage)
	case types.ClientSideStreamingType, types.BidirectionalStreamingType:
		goPackages = strutils.UniqueAndSorted(method.OutType.GoPackage, service.GoPackage)
	}

	data := SubstitutionMethodForStubsView{
		GoPackages:       goPackages,
		Method:           method.Name,
		Service:          service.Name,
		MethodInName:     method.InType.Name,
		MethodOutName:    method.OutType.Name,
		MethodPackage:    strutils.ToSnakeCase(method.Package),
		PackageHeader:    headerPackage(service.GoPackage, service.Name),
		Package:          strutils.ToPackageName(service.GoPackage),
		MethodInPackage:  strutils.ToPackageName(method.InType.GoPackage),
		MethodOutPackage: strutils.ToPackageName(method.OutType.GoPackage),
		URL:              fmt.Sprintf("%s/%s/%s", baseURL, service.Name, method.Name),
	}

	return data, nil
}

func substitutionServiceForStubs(goPackage, serviceName string) SubstitutionServiceForStubsView {
	data := SubstitutionServiceForStubsView{
		GoPackage:     goPackage,
		Service:       serviceName,
		Package:       strutils.ToPackageName(goPackage),
		PackageHeader: headerPackage(goPackage, serviceName),
	}

	return data
}

func headerPackage(goPackage, serviceName string) string {
	return strutils.ToSnakeCase(strutils.ToPackageName(goPackage) + serviceName)
}
