package strutils

import (
	"fmt"
	"strings"
)

// Data types are allowed to be used as a code identifiers.
// That's why there keywords only.
var golangKeywords = map[string]struct{}{
	"break": {}, "default": {}, "func": {}, "interface": {}, "select": {},
	"case": {}, "defer": {}, "go": {}, "map": {}, "struct": {}, "chan": {}, "else": {},
	"goto": {}, "package": {}, "switch": {}, "const": {}, "fallthrough": {}, "if": {},
	"continue": {}, "for": {}, "import": {}, "return": {}, "var": {}, "type": {},
}

func ResolveNameIfCollides(keyword string) string {
	if isGolangKeyword(keyword) {
		return resolve(keyword)
	}

	return keyword
}

func isGolangKeyword(keyword string) bool {
	value := strings.ToLower(keyword)
	_, isExist := golangKeywords[value]

	return isExist
}

func resolve(value string) string {
	const suffixToResolveCollision = "resolved"
	return fmt.Sprintf("%s_%s", value, suffixToResolveCollision)
}
