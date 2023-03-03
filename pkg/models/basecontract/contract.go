package basecontract

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

type Contract struct {
	HeaderPath string

	ImportsPaths []string

	ContractHasMethods bool

	descriptor any
}

func (c Contract) Descriptor() any {
	return c.descriptor
}

func ProtoFromAny(contract Contract) (*desc.FileDescriptor, error) {
	descriptor, ok := contract.Descriptor().(*desc.FileDescriptor)
	if !ok {
		return nil, fmt.Errorf("no valid descriptor")
	}

	return descriptor, nil
}

func OpenAPIFromAny(contract Contract) (*openapi3.T, error) {
	descriptor, ok := contract.Descriptor().(*openapi3.T)
	if !ok {
		return nil, fmt.Errorf("no valid descriptor")
	}

	return descriptor, nil
}

func (c Contract) HasMethods() bool {
	return c.ContractHasMethods
}

func (c Contract) IsProtoContract() bool {
	return fsutils.GetFileExt(c.HeaderPath) == "proto"
}

func (c Contract) IsOpenAPIContract() bool {
	return fsutils.GetFileExt(c.HeaderPath) == "yaml"
}

// SetOfContracts is just alias for several contracts.
type SetOfContracts []Contract

// WithMethodsOnly returns only contracts with methods.
func (set SetOfContracts) WithMethodsOnly() SetOfContracts {
	var filtered SetOfContracts

	for _, contract := range set {
		if contract.HasMethods() {
			filtered = append(filtered, contract)
		}
	}

	return filtered
}

func (set SetOfContracts) HasContractsWithMethods() bool {
	for _, contract := range set {
		if contract.HasMethods() {
			return true
		}
	}

	return false
}

func (set SetOfContracts) FileType() types.SourceFileType {
	switch {
	case sliceutils.FirstOf(set).IsProtoContract():
		return types.ProtoType
	case sliceutils.FirstOf(set).IsOpenAPIContract():
		return types.OpenAPIType
	default:
		return types.UnknownFileType
	}
}
