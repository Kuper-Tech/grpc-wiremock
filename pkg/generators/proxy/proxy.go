package proxy

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
)

// errHasNoContractsWithMethods indicates that provided set of contracts has no methods to mock.
var errHasNoContractsWithMethods = fmt.Errorf("provided contracts has no methods")

// errHasNoProtoContracts indicates that provided set of contracts is not supported by proxy generator.
var errHasNoProtoContracts = fmt.Errorf("provided contracts are not supported")

// compiler abstracts how exactly project should be compiled.
type compiler interface {
	CompileToGo(context.Context, protocontract.Contract, string, io.Writer) error
}

// renderer abstracts how exactly project should be rendered.
type renderer interface {
	Substitute(string, interface{}) (string, error)
}

type proxyGenerator struct {
	port   string
	host   string
	output string

	fs afero.Fs

	renderer
	compiler
	io.Writer
}

func NewGenerator(fs afero.Fs, compiler compiler, renderer renderer, baseURL, output string, logs io.Writer) (proxyGenerator, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return proxyGenerator{}, fmt.Errorf("invalid base url: %w", err)
	}

	return proxyGenerator{
		fs: fs, port: parsed.Port(), host: fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host),
		output: output, renderer: renderer, compiler: compiler, Writer: logs,
	}, nil
}

// Generate generates project template, grpc packages and stubs.
// The result of generation is ready to 'go run' proxy project.
func (g proxyGenerator) Generate(ctx context.Context, contracts protocontract.SetOfContracts) error {
	if !contracts.HasContractsWithMethods() {
		return errHasNoContractsWithMethods
	}

	if !contracts.HasContractsWithProto() {
		return errHasNoProtoContracts
	}

	if err := g.GenerateProject(contracts); err != nil {
		return fmt.Errorf("generate project: %w", err)
	}

	if err := g.GeneratePackages(ctx, contracts, g.Writer); err != nil {
		return fmt.Errorf("generate packages: %w", err)
	}

	if err := g.GenerateStubs(contracts); err != nil {
		return fmt.Errorf("generate stubs: %w", err)
	}

	return nil
}
