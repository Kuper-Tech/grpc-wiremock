package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/generators/certificates/test"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configopener"
)

func Test_testCertificates(t *testing.T) {
	t.Run("with correct certificates", func(t *testing.T) {
		opener := configopener.New(staticFS, "tests/data/wiremock/configs/one-service")
		g := NewGenerator(staticFS, opener)

		caCertPair, certPair, err := g.generate([]string{"awesome"}, "")
		require.NoError(t, err)

		serverConfig, clientConfig, err := createTLSConfig(caCertPair, certPair)
		require.NoError(t, err)

		err = runTest(test.NewCertsTester(serverConfig, clientConfig))
		require.NoError(t, err)
	})
}

func runTest(t test.CertsTester) error {
	t.Server.StartTLS()
	defer t.Server.Close()

	if err := t.DoTestRequest(); err != nil {
		return err
	}

	return nil
}

func createTLSConfig(caPair, certificatePair pair) (*tls.Config, *tls.Config, error) {
	tlsCertificate, err := tls.X509KeyPair(certificatePair.certificate, certificatePair.privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("tls pair: %w", err)
	}

	pool := x509.NewCertPool()

	if !pool.AppendCertsFromPEM(caPair.certificate) {
		return nil, nil, fmt.Errorf("pool doesn't contain certificate authority")
	}

	clientTLSConf := tls.Config{RootCAs: pool}
	serverTLSConf := tls.Config{Certificates: []tls.Certificate{tlsCertificate}}

	return &serverTLSConf, &clientTLSConf, nil
}
