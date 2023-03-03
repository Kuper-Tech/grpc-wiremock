package proxy

import (
	"fmt"
	"path/filepath"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/protocontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

const (
	mainGoTemplatePath = "proxy/files/main.go.tpl"
	projectTemplateDir = "proxy/template/layout"
)

// GenerateProject generates golang project based on the template.
// And calls proto compiler to compile packages based on provided proto contracts.
func (g proxyGenerator) GenerateProject(contracts protocontract.SetOfContracts) error {
	return g.generateProject(contracts, projectTemplateDir)
}

func (g proxyGenerator) generateProject(contracts protocontract.SetOfContracts, templatePath string) error {
	staticFS := static.FromEmbed()

	if err := fsutils.CopyDir(staticFS, g.fs, templatePath, g.output, false); err != nil {
		return fmt.Errorf("copy template: %w", err)
	}

	fileToPath := map[string]struct {
		path         string
		substitution interface{}
	}{
		"cmd/main.go": {
			path:         mainGoTemplatePath,
			substitution: substitutionForProject(contracts),
		},
	}

	for outputFilePath, template := range fileToPath {
		content, err := g.renderer.Substitute(template.path, &template.substitution)
		if err != nil {
			return fmt.Errorf("substitute file: %s, err: %w", template.path, err)
		}

		pathInProject := filepath.Join(g.output, outputFilePath)

		if err = fsutils.WriteFile(g.fs, pathInProject, content); err != nil {
			return fmt.Errorf("write main.go file: %w", err)
		}
	}

	return nil
}
