package protocontract

import (
	"fmt"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/blacklist"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

// FromProtoDescriptor converts proto descriptor into Contract model.
func FromProtoDescriptor(descriptor *desc.FileDescriptor, protoPath string) (Contract, error) {
	goPackage := protoGoPackage(descriptor)

	services, err := protoServices(descriptor)
	if err != nil {
		return Contract{}, fmt.Errorf("get services: %w", err)
	}

	base := Base{
		Services:     services,
		GoPackage:    goPackage,
		Package:      descriptor.GetPackage(),
		Messages:     protoMessages(descriptor),
		ImportsPaths: protoImports(descriptor, protoPath),
		HeaderPath:   filepath.Join(protoPath, descriptor.GetName()),
	}

	contract := Contract{
		Base: base,
		Desc: descriptor,
	}

	return contract, nil
}

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

func protoImports(descriptor *desc.FileDescriptor, protoPath string) []string {
	var imports []string

	for _, path := range getImports(descriptor) {
		if blacklist.IsDeliveredWithProtoc(path) {
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

func protoServices(descriptor *desc.FileDescriptor) ([]Service, error) {
	var services []Service

	for _, service := range descriptor.GetServices() {
		var methods []Method

		for _, method := range service.GetMethods() {
			in, err := protoMessageType(method, method.GetInputType())
			if err != nil {
				return nil, fmt.Errorf("message-in type: %w", err)
			}

			out, err := protoMessageType(method, method.GetOutputType())
			if err != nil {
				return nil, fmt.Errorf("message-out type: %w", err)
			}

			methods = append(methods, Method{
				InType:     in,
				OutType:    out,
				Name:       method.GetName(),
				Package:    method.GetFile().GetPackage(),
				MethodType: types.MethodType(method.IsClientStreaming(), method.IsServerStreaming()),
			})
		}

		services = append(services, Service{
			Methods:   methods,
			Name:      service.GetName(),
			GoPackage: protoGoPackage(descriptor),
		})
	}

	return services, nil
}

func protoMessages(descriptor *desc.FileDescriptor) []string {
	var messages []string

	for _, message := range descriptor.GetMessageTypes() {
		messages = append(messages, message.GetName())
	}

	return messages
}

func protoMessageType(methodDesc *desc.MethodDescriptor, messageDesc *desc.MessageDescriptor) (MessageType, error) {
	goPackage := protoGoPackage(messageDesc.GetFile())

	message := MessageType{
		GoPackage:  goPackage,
		Name:       messageDesc.GetName(),
		Package:    messageDesc.GetFile().GetPackage(),
		IsExternal: messageDesc.GetFile().GetPackage() != methodDesc.GetFile().GetPackage(),
	}

	return message, nil
}
