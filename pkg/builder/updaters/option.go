package updaters

import (
	"fmt"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/genproto/googleapis/api/annotations"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

type optionUpdater struct {
	annotationsDescriptor *desc.FileDescriptor
}

func NewOptionUpdater() (*optionUpdater, error) {
	descriptor, err := parseAnnotations()
	if err != nil {
		return nil, fmt.Errorf("get annotations descriptor: %w", err)
	}

	return &optionUpdater{annotationsDescriptor: descriptor}, nil
}

func (u *optionUpdater) Name() string {
	return "option updater"
}

func (u *optionUpdater) Update(contract *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	return u.overwriteOptions(contract)
}

func (u *optionUpdater) overwriteOptions(contract *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	withAnnotations, err := includeAnnotations(contract, u.annotationsDescriptor)
	if err != nil {
		return nil, fmt.Errorf("include annotations: %w", err)
	}

	withOptions, err := addOptions(withAnnotations)
	if err != nil {
		return nil, fmt.Errorf("add options: %w", err)
	}

	return withOptions, nil
}

func addOptions(descriptor *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	fd := *descriptor

	for _, service := range fd.GetServices() {
		if err := addOptionsForService(service); err != nil {
			return nil, fmt.Errorf("for service: %w", err)
		}
	}

	return &fd, nil
}

func addOptionsForService(descriptor *desc.ServiceDescriptor) error {
	methods := descriptor.GetMethods()

	for idx, method := range methods {
		url, err := createRpcUrl(descriptor.GetName(), method.GetName())
		if err != nil {
			return fmt.Errorf("create rpc url: %w", err)
		}

		updatedDescriptor, err := addOptionForMethod(method, url)
		if err != nil {
			return fmt.Errorf("update method descriptor: %w", err)
		}

		methods[idx] = updatedDescriptor
	}

	return nil
}

func addOptionForMethod(descriptor *desc.MethodDescriptor, url string) (*desc.MethodDescriptor, error) {
	methodBuilder, err := builder.FromMethod(descriptor)
	if err != nil {
		return nil, fmt.Errorf("create method builder: %w", err)
	}

	opt, err := createOption(url)
	if err != nil {
		return nil, fmt.Errorf("create method option: %w", err)
	}

	updatedDescriptor, err := methodBuilder.SetOptions(opt).Build()
	if err != nil {
		return nil, fmt.Errorf("set option: %w", err)
	}

	return updatedDescriptor, nil
}

func parseAnnotations() (*desc.FileDescriptor, error) {
	parser := &protoparse.Parser{
		IncludeSourceCodeInfo: true,

		ImportPaths: []string{environment.TmpAnnotationProtosDir},
	}

	annotationFiles := []string{
		environment.AnnotationsPath,
		environment.AnnotationsHttpPath,
	}

	annotationDesc, err := parser.ParseFiles(annotationFiles...)
	if err != nil {
		return nil, fmt.Errorf("parse 'annotations.proto': %w", err)
	}

	return sliceutils.FirstOf(annotationDesc), nil
}

func includeAnnotations(contractDescriptor, annotationsDescriptor *desc.FileDescriptor) (*desc.FileDescriptor, error) {
	fileBuilder, err := builder.FromFile(contractDescriptor)
	if err != nil {
		return nil, fmt.Errorf("create file builder: %w", err)
	}

	updatedFileDescriptor, err := fileBuilder.AddImportedDependency(annotationsDescriptor).Build()
	if err != nil {
		return nil, fmt.Errorf("add dependency: %w", err)
	}

	return updatedFileDescriptor, nil
}

func createOption(mockURL string) (*dpb.MethodOptions, error) {
	//option (google.api.http) = {
	//	post: "/v1/example/echo"
	//	body: "*"
	//};

	options := dpb.MethodOptions{}
	httpRule := annotations.HttpRule{
		Body:    "*",
		Pattern: &annotations.HttpRule_Post{Post: mockURL},
	}

	if err := proto.SetExtension(&options, annotations.E_Http, &httpRule); err != nil {
		return nil, fmt.Errorf("set extension: %w", err)
	}

	return &options, nil
}

func createRpcUrl(serviceName string, methodName string) (string, error) {
	if serviceName == "" || methodName == "" {
		return "", fmt.Errorf("empty name")
	}

	return filepath.Join(
		"/",
		strutils.ToCamelCase(serviceName),
		strutils.ToCamelCase(methodName),
	), nil
}
