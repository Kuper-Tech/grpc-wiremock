package certgen

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/certificates"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configopener"
)

type certsGen struct {
	supervisordPath string

	output string

	logger io.Writer

	afero.Fs
}

func NewCertsGenWithDefaultFs(output, supervisordPath string, logger io.Writer) certsGen {
	return certsGen{output: output, supervisordPath: supervisordPath, logger: logger, Fs: afero.NewOsFs()}
}

func (g *certsGen) Generate(ctx context.Context) error {
	opener := configopener.New(g.Fs, g.supervisordPath)
	generator := certificates.NewGenerator(g.Fs, opener)

	if err := generator.Generate(ctx, g.output); err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	return nil
}
