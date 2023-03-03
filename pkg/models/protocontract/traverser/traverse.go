package traverser

import (
	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
)

func Descriptors(contracts protocontract.SetOfContracts) []*desc.FileDescriptor {
	var all []*desc.FileDescriptor

	for _, cont := range contracts {
		all = append(all, cont.Desc)
		rec(cont.Desc, &all)
	}

	return all
}

func rec(desc *desc.FileDescriptor, all *[]*desc.FileDescriptor) {
	*all = append(*all, desc.GetDependencies()...)
	for _, descriptor := range desc.GetDependencies() {
		rec(descriptor, all)
	}
}
