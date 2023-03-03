package basecontract

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func FromOpenAPIDescriptor(descriptor *openapi3.T, path string) (Contract, error) {
	if descriptor == nil {
		return Contract{}, fmt.Errorf("descriptor is nil")
	}

	contract := Contract{
		HeaderPath: path,

		ContractHasMethods: len(descriptor.Paths) != 0,

		descriptor: descriptor,
	}

	return contract, nil
}
