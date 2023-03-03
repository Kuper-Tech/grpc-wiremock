package certificates

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"path/filepath"
	"time"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

const (
	defaultBitsCount = 2048

	serialNumber = 1126

	privateKeyPEMBlockName  = "RSA PRIVATE KEY"
	certificatePEMBlockName = "CERTIFICATE"
)

type configOpener interface {
	Open() (config.Wiremock, error)
}

var caCertificateDesc = &x509.Certificate{
	IsCA: true,

	Subject: defaultSubject,

	SerialNumber: big.NewInt(serialNumber),

	NotBefore: time.Now(),
	NotAfter:  todayPlusOneYear(),

	KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},

	BasicConstraintsValid: true,
}

var defaultSubject = pkix.Name{
	CommonName: "grpc-wiremock",

	Country:  []string{"RU"},
	Locality: []string{"Moscow"},

	Organization: []string{"Sbermarket LLC"},
}

type pair struct {
	privateKey  []byte
	certificate []byte
}

type certsGen struct {
	opener configOpener

	fs afero.Fs
}

func NewGenerator(fs afero.Fs, opener configOpener) certsGen {
	return certsGen{opener: opener, fs: fs}
}

func (g *certsGen) Generate(_ context.Context, output string) error {
	collectedDomains, err := g.collectDomains(commonDomains, domains)
	if err != nil {
		return fmt.Errorf("generate domains: %w", err)
	}

	caCertPair, certPair, err := g.generate(collectedDomains, output)
	if err != nil {
		return fmt.Errorf("generate cert pairs: %w", err)
	}

	filesToSave := map[string][]byte{
		environment.CertKeyFile:  certPair.privateKey,
		environment.CertCertFile: certPair.certificate,

		environment.CAKeyFile:  caCertPair.privateKey,
		environment.CACertFile: caCertPair.certificate,
	}

	for filePath, content := range filesToSave {
		saveToPath := filepath.Join(output, filePath)

		if err = fsutils.WriteFile(g.fs, saveToPath, string(content)); err != nil {
			return fmt.Errorf("save '%s': %w", saveToPath, err)
		}
	}

	return nil
}

func (g *certsGen) generate(collectedDomains []string, output string) (pair, pair, error) {
	caPrivateKey, err := g.getCAKey(output)
	if err != nil {
		return pair{}, pair{}, fmt.Errorf("get key: %w", err)
	}

	caPair, err := g.getCertificateAuthority(caCertificateDesc, caPrivateKey, output)
	if err != nil {
		return pair{}, pair{}, fmt.Errorf("generate ca: %w", err)
	}

	certificateDesc := createCertificateDesc(collectedDomains)

	certificatePair, err := g.createCertificate(caCertificateDesc, certificateDesc, caPrivateKey)
	if err != nil {
		return pair{}, pair{}, fmt.Errorf("generate certificate: %w", err)
	}

	return caPair, certificatePair, nil
}

func (g *certsGen) getCAKey(output string) (*rsa.PrivateKey, error) {
	var caPrivateKey *rsa.PrivateKey
	var err error

	if !environment.IsCAExists(g.fs, output) {
		caPrivateKey, err = rsa.GenerateKey(rand.Reader, defaultBitsCount)
		if err != nil {
			return nil, fmt.Errorf("create ca private key: %w", err)
		}

		return caPrivateKey, nil
	}

	keyFilePath := filepath.Join(output, environment.CAKeyFile)

	keyContent, err := afero.ReadFile(g.fs, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	keyPEM, err := unmarshalFromPEMBlock(keyContent)
	if err != nil {
		return nil, fmt.Errorf("unmarshal key file: %w", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse key file: %w", err)
	}

	return privateKey, nil
}

func createCertificateDesc(dnsNames []string) *x509.Certificate {
	return &x509.Certificate{
		Subject: defaultSubject,

		DNSNames: dnsNames,
		KeyUsage: x509.KeyUsageDigitalSignature,

		SerialNumber: big.NewInt(serialNumber),

		NotBefore: time.Now(),
		NotAfter:  todayPlusOneYear(),

		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}
}

func marshalToPEMBlocks(contentType string, content []byte) ([]byte, error) {
	var buffer bytes.Buffer

	block := pem.Block{Type: contentType, Bytes: content}

	if err := pem.Encode(&buffer, &block); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return buffer.Bytes(), nil
}

func unmarshalFromPEMBlock(content []byte) (*pem.Block, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("empty block")
	}

	return block, nil
}

func todayPlusOneYear() time.Time {
	return time.Now().AddDate(1, 0, 0)
}
