package renderer

import (
	"path"
	"strings"
	"text/template"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
)

var funcMap = template.FuncMap{
	"Base":          path.Base,
	"Capitalize":    strings.Title,
	"ToLower":       strings.ToLower,
	"ToSnakeCase":   strutils.ToSnakeCase,
	"ToPackageName": strutils.ToPackageName,
}
