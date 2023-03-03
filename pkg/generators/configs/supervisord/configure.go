package supervisord

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/configs"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

// renderer abstracts how exactly project should be rendered.
type renderer interface {
	Substitute(string, interface{}) (string, error)
}

type supervisordDomainConfView struct {
	Domain string
	Port   string
	Root   string
}

type Configuer struct {
	Renderer   renderer
	OutputPath string

	afero.Fs
}

func (c Configuer) GenerateConfig(values configs.Values) error {
	const templatePath = "supervisord/files/supervisord.conf.tpl"

	confView := supervisordDomainConfView{Domain: values.Domain, Port: values.Port, Root: values.Root}

	content, err := c.Renderer.Substitute(templatePath, &confView)
	if err != nil {
		return fmt.Errorf("substitute: %w", err)
	}

	pathToSave := c.createConfigPath(values.Domain)

	if err = fsutils.WriteFile(c.Fs, pathToSave, content); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (c Configuer) createConfigPath(domain string) string {
	return filepath.Join(c.OutputPath, fmt.Sprintf("mock-%s.conf", domain))
}
