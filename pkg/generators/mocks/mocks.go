package mocks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/jsonschema/schematomock"
	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/mock"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/client"
)

// errHasNoContractsWithMethods indicates that provided set of contracts has no methods to mock.
var errHasNoContractsWithMethods = fmt.Errorf("provided contracts has no methods")

type mocksGenerator struct {
	fs afero.Fs

	logs io.Writer

	outputPath string
}

func NewGenerator(fs afero.Fs, outputPath string, logs io.Writer) mocksGenerator {
	return mocksGenerator{fs: fs, logs: logs, outputPath: outputPath}
}

func (g mocksGenerator) Generate(ctx context.Context, contracts contract.SetOfContracts, domain string) error {
	if !contracts.HasContractsWithMethods() {
		return errHasNoContractsWithMethods
	}

	for _, contract := range contracts {
		if !contract.IsOpenAPIContract() {
			continue
		}

		if err := g.generate(ctx, contract, domain); err != nil {
			return fmt.Errorf("generate mocks for contract '%s': %w", contract.HeaderPath, err)
		}
	}

	return nil
}

func (g mocksGenerator) generate(_ context.Context, contract contract.Contract, domain string) error {
	descriptor, ok := contract.Descriptor().(*openapi3.T)
	if !ok {
		return fmt.Errorf("cast openapi desc")
	}

	mocks, err := schematomock.Parse(descriptor)
	if err != nil {
		return fmt.Errorf("load descriptor: %s", contract.HeaderPath)
	}

	if err = g.generateWiremockMocks(domain, mocks); err != nil {
		return fmt.Errorf("generate mocks: %s, err: %w", contract.HeaderPath, err)
	}

	return nil
}

func (g mocksGenerator) generateWiremockMocks(domain string, values []mock.Mock) error {
	const mappingsDir = "mappings"

	for _, value := range values {
		mockName, mockContent, err := generateWiremockMock(value)
		if err != nil {
			return fmt.Errorf("generate mock: %w", err)
		}

		mockPath := filepath.Join(g.outputPath, domain, mappingsDir, mockName)

		if err = fsutils.WriteFile(g.fs, mockPath, mockContent); err != nil {
			return fmt.Errorf("write mock: %w", err)
		}
	}

	return nil
}

func generateWiremockMock(mock mock.Mock) (string, string, error) {
	name := createMockName(mock)

	wiremockMock := client.DefaultMock()

	wiremockMock.WithRequestMethod(mock.RequestMethod).
		WithRequestUrlPath(mock.RequestUrlPath).
		WithResponseStatusCode(mock.ResponseStatusCode).
		WithResponseBody(mock.ResponseBody).
		WithName(strings.TrimSuffix(name, ".json")).
		WithDescription(mock.Description)

	content, err := json.MarshalIndent(wiremockMock, "", "\t")
	if err != nil {
		return "", "", fmt.Errorf("marshal mock content: %w", err)
	}

	return name, string(content), nil
}

func createMockName(mock mock.Mock) string {
	var parts []string

	shouldSkipPart := func(part string) bool {
		const space = " "
		return len(part) == 0 || part == space
	}

	split := strings.Split(mock.RequestUrlPath, "/")
	for _, part := range split {
		if shouldSkipPart(part) {
			continue
		}
		parts = append(parts, strutils.ToSnakeCase(part))
	}

	return fmt.Sprintf(
		"%s_%s_%d.json",
		strutils.ToSnakeCase(mock.RequestMethod),
		strings.Join(parts, "_"),
		mock.ResponseStatusCode,
	)
}
