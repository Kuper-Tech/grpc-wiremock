package strutils

import (
	"bytes"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")
	output = strings.ToLower(output)
	return strings.ReplaceAll(output, ".", "_")
}

func ToCamelCase(input string) string {
	res := bytes.NewBuffer(nil)
	capNext := true

	for _, v := range input {
		if unicode.IsUpper(v) {
			res.WriteRune(v)
			capNext = false
			continue
		}

		if unicode.IsDigit(v) {
			res.WriteRune(v)
			capNext = true
			continue
		}

		if unicode.IsLower(v) {
			if capNext {
				res.WriteRune(unicode.ToUpper(v))
			} else {
				res.WriteRune(v)
			}
			capNext = false
			continue
		}

		capNext = true
	}

	return res.String()
}

func ToPackageName(input string) string {
	return ToSnakeCase(filepath.Base(input))
}

func UniqueAndSorted(values ...string) []string {
	uniq := map[string]struct{}{}
	for _, value := range values {
		uniq[value] = struct{}{}
	}

	var uniqValues []string
	for value := range uniq {
		uniqValues = append(uniqValues, value)
	}

	sort.Slice(uniqValues, func(i, j int) bool {
		return uniqValues[i] < uniqValues[j]
	})

	return uniqValues
}
