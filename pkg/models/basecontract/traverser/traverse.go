package traverser

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
)

func Descriptors(contracts contract.SetOfContracts) ([]*desc.FileDescriptor, error) {
	var all []*desc.FileDescriptor

	for _, cont := range contracts {
		descriptor, err := contract.ProtoFromAny(cont)
		if err != nil {
			return nil, fmt.Errorf("from any: %w", err)
		}

		all = append(all, descriptor)
		rec(descriptor, &all)
	}

	return all, nil
}

func rec(desc *desc.FileDescriptor, all *[]*desc.FileDescriptor) {
	*all = append(*all, desc.GetDependencies()...)
	for _, descriptor := range desc.GetDependencies() {
		rec(descriptor, all)
	}
}
