package proxy

import (
	"fmt"
	"strings"
)

const pkgDir = "grpc-proxy/pkg"

// toOriginalPackageName helps to convert modified for
// generation purposes `go_package` to original one.
func toOriginalPackageName(goPackage string) string {
	return strings.TrimPrefix(goPackage, fmt.Sprintf("%s/", pkgDir))
}
