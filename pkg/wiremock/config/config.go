package config

import (
	"fmt"
	"path/filepath"
)

var (
	EmptyWiremockConfigErr = fmt.Errorf("empty wiremock config")

	NoCorrespondingAPIErr = fmt.Errorf("no API corresponding to the provided domain")
)

type Service struct {
	Name    string `json:"name"`
	Port    int    `json:"port"`
	RootDir string `json:"rootDir"`
}

type Wiremock struct {
	Services []Service `json:"services"`
}

func NewService(root, domain string, port int) Service {
	return Service{Name: domain, Port: port, RootDir: filepath.Join(root, domain)}
}
