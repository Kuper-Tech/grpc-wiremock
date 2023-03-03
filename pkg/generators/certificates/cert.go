package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
)

func (g *certsGen) getCertificateAuthority(caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, output string) (pair, error) {
	var caContent []byte
	var err error

	convertedPrivateKey := x509.MarshalPKCS1PrivateKey(caPrivateKey)

	caPrivateKeyPEM, err := marshalToPEMBlocks(privateKeyPEMBlockName, convertedPrivateKey)
	if err != nil {
		return pair{}, fmt.Errorf("marshal ca private key: %w", err)
	}

	if environment.IsCAExists(g.fs, output) {
		caCertFilePath := filepath.Join(output, environment.CACertFile)

		caPEM, err := afero.ReadFile(g.fs, caCertFilePath)
		if err != nil {
			return pair{}, fmt.Errorf("read ca content: %w", err)
		}

		return pair{privateKey: caPrivateKeyPEM, certificate: caPEM}, nil
	}

	caContent, err = x509.CreateCertificate(rand.Reader, caCertificate, caCertificate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return pair{}, fmt.Errorf("create ca certificate: %w", err)
	}

	caPEM, err := marshalToPEMBlocks(certificatePEMBlockName, caContent)
	if err != nil {
		return pair{}, fmt.Errorf("marshal ca: %w", err)
	}

	return pair{privateKey: caPrivateKeyPEM, certificate: caPEM}, nil
}

func (g *certsGen) createCertificate(caCertificateDesc, certificateDesc *x509.Certificate, caPrivateKey *rsa.PrivateKey) (pair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, defaultBitsCount)
	if err != nil {
		return pair{}, fmt.Errorf("generate private key: %w", err)
	}

	certificateContent, err := x509.CreateCertificate(rand.Reader, certificateDesc, caCertificateDesc, &privateKey.PublicKey, caPrivateKey)
	if err != nil {
		return pair{}, fmt.Errorf("create certificate: %w", err)
	}

	convertedPrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM, err := marshalToPEMBlocks(privateKeyPEMBlockName, convertedPrivateKey)
	if err != nil {
		return pair{}, fmt.Errorf("marshal private key: %w", err)
	}

	certificatePEM, err := marshalToPEMBlocks(certificatePEMBlockName, certificateContent)
	if err != nil {
		return pair{}, fmt.Errorf("marshal certificate: %w", err)
	}

	return pair{privateKey: privateKeyPEM, certificate: certificatePEM}, nil
}
