package client

import (
	"net/http"
	"time"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

type configOpener interface {
	Open() (config.Wiremock, error)
}

type wiremock struct {
	host string
	port uint

	client http.Client

	fs afero.Fs

	configOpener
}

func NewDefaultClient(fs afero.Fs, opener configOpener) *wiremock {
	var defaultTimeout = 15 * time.Second

	const (
		defaultHost = "http://localhost"
		defaultPort = 9000
	)

	return NewClient(fs, opener, defaultHost, defaultPort, http.Client{Timeout: defaultTimeout})
}

func NewClient(fs afero.Fs, opener configOpener, host string, port uint, httpClient http.Client) *wiremock {
	return &wiremock{fs: fs, configOpener: opener, host: host, port: port, client: httpClient}
}
