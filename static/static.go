package static

import (
	"embed"

	"github.com/spf13/afero"
)

//go:embed *
var static embed.FS

func FromEmbed() afero.Fs {
	return afero.FromIOFS{FS: static}
}
