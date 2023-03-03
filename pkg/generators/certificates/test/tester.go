package test

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

const successMessage = "success"

type CertsTester struct {
	Client *http.Client
	Server *httptest.Server
}

func NewCertsTester(serverCert, clientCert *tls.Config) CertsTester {
	server := httptest.NewUnstartedServer(createTestHandler())
	server.TLS = serverCert

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientCert,
		},
	}

	return CertsTester{Server: server, Client: client}
}

func (t *CertsTester) DoTestRequest() error {
	response, err := t.Client.Get(t.Server.URL)
	if err != nil {
		return fmt.Errorf("do get: %w", err)
	}

	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	return requireEqualBody(string(responseContent), successMessage)
}

func createTestHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, successMessage)
	}
}

func requireEqualBody(got, target string) error {
	if strings.TrimSpace(got) != target {
		return fmt.Errorf("not equal got: %s, target: %s", got, target)
	}

	return nil
}
