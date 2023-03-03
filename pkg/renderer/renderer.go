package renderer

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
)

type renderer struct {
	fs afero.Fs

	tpls map[string]*template.Template

	pathToTemplates string
}

func New(fs afero.Fs, path string) (renderer, error) {
	newRenderer := renderer{
		fs:              fs,
		pathToTemplates: path,

		tpls: map[string]*template.Template{},
	}

	templatesGlob := filepath.Join(path, "*")

	files, err := afero.Glob(fs, templatesGlob)
	if err != nil {
		return renderer{}, fmt.Errorf("glob: %w", err)
	}

	for _, name := range files {
		data, err := afero.ReadFile(fs, name)
		if err != nil {
			return renderer{}, fmt.Errorf("ReadFile: %w", err)
		}

		tpl, err := template.New(name).Funcs(funcMap).Parse(string(data))
		if err != nil {
			return renderer{}, fmt.Errorf("parse: %w", err)
		}

		newRenderer.tpls[name] = tpl
	}

	return newRenderer, nil
}

// Substitute tries to correct collisions with go keywords. The data should be passed as a pointer.
// The data will then be substituted into the template from static FS.
func (s renderer) Substitute(name string, v interface{}) (string, error) {
	resolveCollisions(v)

	var buff strings.Builder
	tpl, ok := s.tpls[name]
	if !ok {
		return "", fmt.Errorf("template not found: %s", name)
	}

	if err := tpl.Execute(&buff, v); err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}

	return buff.String(), nil
}
