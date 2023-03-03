package basecontract

import (
	"fmt"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/blacklist"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
)

// FromProtoDescriptor converts proto descriptor into Contract model.
func FromProtoDescriptor(descriptor *desc.FileDescriptor, protoPath string) (Contract, error) {
	if descriptor == nil {
		return Contract{}, fmt.Errorf("descriptor is nil")
	}

	contract := Contract{
		descriptor:         descriptor,
		ContractHasMethods: hasMethods(descriptor),
		HeaderPath:         filepath.Join(protoPath, descriptor.GetName()),

		ImportsPaths: protoImports(descriptor, protoPath),
	}

	return contract, nil
}

func protoImports(descriptor *desc.FileDescriptor, protoPath string) []string {
	var imports []string

	for _, path := range getImports(descriptor) {
		if blacklist.IsDeliveredWithProtoc(path) {
			continue
		}

		if blacklist.IsGoogleAPIContract(path) {
			continue
		}

		imports = append(imports, filepath.Join(protoPath, path))
	}

	return imports
}

func getImports(descriptor *desc.FileDescriptor) []string {
	var all []string
	traverse(descriptor, &all)

	var paths []string
	for _, path := range all {
		paths = append(paths, path)
	}

	return paths
}

func traverse(desc *desc.FileDescriptor, all *[]string) {
	for _, descriptor := range desc.GetDependencies() {
		*all = append(*all, descriptor.GetName())
		traverse(descriptor, all)
	}
}

func hasMethods(descriptor *desc.FileDescriptor) bool {
	for _, service := range descriptor.GetServices() {
		if len(service.GetMethods()) > 0 {
			return true
		}
	}

	return false
}
