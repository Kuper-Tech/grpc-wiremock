package blacklist

import "strings"

func IsDeliveredWithProtoc(value string) bool {
	var (
		filePathsDoNotTouch  = []string{"google/protobuf"}
		goPackagesDoNotTouch = []string{"google.golang.org/protobuf"}
	)

	for _, file := range filePathsDoNotTouch {
		if strings.Contains(value, file) {
			return true
		}
	}

	for _, goPackage := range goPackagesDoNotTouch {
		if strings.Contains(value, goPackage) {
			return true
		}
	}

	return false
}

func IsGoogleAPIContract(value string) bool {
	var (
		filePathsDoNotTouch  = []string{"google/api"}
		goPackagesDoNotTouch = []string{"google.golang.org/genproto/googleapis/api"}
	)

	for _, file := range filePathsDoNotTouch {
		if strings.Contains(value, file) {
			return true
		}
	}

	for _, goPackage := range goPackagesDoNotTouch {
		if strings.Contains(value, goPackage) {
			return true
		}
	}

	return false
}
