package builder

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/blacklist"
)

type protoUpdater interface {
	Name() string
	Update(contract *desc.FileDescriptor) (*desc.FileDescriptor, error)
}

func UpdateContracts(contracts []*desc.FileDescriptor, updaters ...protoUpdater) ([]*desc.FileDescriptor, error) {
	var updatedDescriptors []*desc.FileDescriptor

	for _, descriptor := range contracts {
		goPackage := protoGoPackage(descriptor)

		if blacklist.IsGoogleAPIContract(goPackage) {
			updatedDescriptors = append(updatedDescriptors, descriptor)
			continue
		}

		if blacklist.IsDeliveredWithProtoc(goPackage) {
			updatedDescriptors = append(updatedDescriptors, descriptor)
			continue
		}

		updatedDescriptor := descriptor

		var err error
		for _, updater := range updaters {
			updatedDescriptor, err = updater.Update(updatedDescriptor)
			if err != nil {
				return nil, fmt.Errorf("update proto descriptor with %s: %w", updater.Name(), err)
			}
		}

		updatedDescriptors = append(updatedDescriptors, updatedDescriptor)
	}

	return updatedDescriptors, nil
}

// TODO separate common functions
func protoGoPackage(descriptor *desc.FileDescriptor) string {
	if descriptor.GetFileOptions() == nil {
		return ""
	}

	goPackage := descriptor.GetFileOptions().GoPackage

	if goPackage == nil {
		return ""
	}

	return *goPackage
}
