package proxy

import (
	"fmt"
	"path/filepath"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

const (
	serviceTemplateName = "proxy/files/service.go.tpl"

	unaryMethodTemplateName = "proxy/files/method-unary.go.tpl"

	clientStreamingMethodTemplateName = "proxy/files/method-client-streaming.go.tpl"

	serverStreamingMethodTemplateName = "proxy/files/method-server-streaming.go.tpl"

	bidirectionalStreamingMethodTemplateName = "proxy/files/method-bidirectional-streaming.go.tpl"
)

// GenerateStubs generates golang grpc stubs based on the template and the provided contracts.
func (g proxyGenerator) GenerateStubs(contracts protocontract.SetOfContracts) error {
	return g.generateStubs(contracts)
}

func (g proxyGenerator) generateStubs(contracts protocontract.SetOfContracts) error {
	for _, contractToGenerate := range contracts {
		if err := g.generateForContract(contractToGenerate); err != nil {
			return fmt.Errorf("generate stubs for contract '%s', err: %w", contractToGenerate.HeaderPath, err)
		}
	}

	return nil
}

func (g proxyGenerator) generateForContract(contractToGenerate protocontract.Contract) error {
	for _, service := range contractToGenerate.Services {
		if err := g.generateForService(contractToGenerate, service); err != nil {
			return fmt.Errorf("generate for service '%s', err: %w", service.Name, err)
		}
	}

	return nil
}

func (g proxyGenerator) generateForService(contractToGenerate protocontract.Contract, service protocontract.Service) error {
	for _, method := range service.Methods {
		if err := g.generateForMethod(service, method); err != nil {
			return fmt.Errorf("generate for method '%s', err: %w", method.Name, err)
		}
	}

	templatePath := serviceTemplateName
	substitution := substitutionServiceForStubs(contractToGenerate.GoPackage, service.Name)

	content, err := g.renderer.Substitute(templatePath, &substitution)
	if err != nil {
		return fmt.Errorf("substitute file: %s, err: %w", templatePath, err)
	}

	pathInProject := filepath.Join(g.output, serviceFileName(service.GoPackage, service.Name))

	if err = fsutils.WriteFile(g.fs, pathInProject, content); err != nil {
		return fmt.Errorf("write service file: %w", err)
	}

	return nil
}

func (g proxyGenerator) generateForMethod(service protocontract.Service, method protocontract.Method) error {
	templateType, err := getMethodTemplate(method.MethodType)
	if err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	substitution, err := substitutionMethodForStubs(service, method, g.host)
	if err != nil {
		return fmt.Errorf("get substitution: %w", err)
	}

	content, err := g.renderer.Substitute(templateType, &substitution)
	if err != nil {
		return fmt.Errorf("substitute file: %s, err: %w", templateType, err)
	}

	pathInProject := filepath.Join(g.output, methodFileName(service.GoPackage, service.Name, method.Name))

	if err = fsutils.WriteFile(g.fs, pathInProject, content); err != nil {
		return fmt.Errorf("write method file: %w", err)
	}

	return nil
}

func getMethodTemplate(methodType types.ProtoMethodType) (string, error) {
	switch methodType {
	case types.UnaryType:
		return unaryMethodTemplateName, nil
	case types.ClientSideStreamingType:
		return clientStreamingMethodTemplateName, nil
	case types.ServerSideStreamingType:
		return serverStreamingMethodTemplateName, nil
	case types.BidirectionalStreamingType:
		return bidirectionalStreamingMethodTemplateName, nil
	}

	return "", fmt.Errorf("incorrect method type %s", methodType)
}

func serviceFileName(packageName, serviceName string) string {
	return filepath.Join("internal", toOriginalPackageName(packageName), strutils.ToSnakeCase(serviceName), "service.go")
}

func methodFileName(packageName, serviceName, methodName string) string {
	return filepath.Join("internal", toOriginalPackageName(packageName), strutils.ToSnakeCase(serviceName), strutils.ToSnakeCase(methodName)+".go")
}
