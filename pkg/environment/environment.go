package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/static"
)

const (
	TmpUnifiedContractsDir     = "/tmp/unified-contracts"
	TmpOverwrittenContractsDir = "/tmp/overwritten-contracts"
	TmpGeneratedPackagesDir    = "/tmp/generated-packages"

	NginxConfigsPath          = "/etc/nginx/http.d"
	SupervisordConfigsPath    = "/etc/supervisord/mocks"
	SupervisordConfigsDirPath = "/etc/supervisord"
	SupervisordMainConfigPath = "/etc/supervisord/supervisord.conf"

	DefaultWiremockConfigPath = "/home/mock"

	TmpWellKnownProtosDir  = "/tmp/proto-includes"
	TmpAnnotationProtosDir = "/tmp/proto-annotations"

	AnnotationsPath     = "google/api/annotations.proto"
	AnnotationsHttpPath = "google/api/http.proto"

	DefaultCertificatesPath = "/etc/ssl"

	TrustedCertificatePath = "/etc/ssl/certs/ca-certificates.crt"

	CAKeyFile  = "mock/share/mockCA.key"
	CACertFile = "mock/share/mockCA.crt"

	CertKeyFile  = "mock/mock.key"
	CertCertFile = "mock/mock.crt"
)

func DumpProtos(fs afero.Fs) error {
	protoToCopy := map[string]string{
		"proto-includes":    TmpWellKnownProtosDir,
		"proto-annotations": TmpAnnotationProtosDir,
	}

	staticFS := static.FromEmbed()

	for sourcePath, targetPath := range protoToCopy {
		if err := fsutils.CopyDir(staticFS, fs, sourcePath, targetPath, true); err != nil {
			return fmt.Errorf("copy protos from embed fs: %w", err)
		}
	}

	return nil
}

func CleanTmpDirs(fs afero.Fs) error {
	tmpDirs := []string{
		TmpUnifiedContractsDir,
		TmpGeneratedPackagesDir,
		TmpOverwrittenContractsDir,
	}

	if err := fsutils.RemoveTmpDirs(fs, tmpDirs...); err != nil {
		return fmt.Errorf("remove tmp dir: %w", err)
	}

	return nil
}

func CleanConfigs(fs afero.Fs, path string) error {
	match := func(info os.FileInfo) bool {
		if info.IsDir() {
			return false
		}

		const (
			mockPrefix = "mock-"
			mockSuffix = ".conf"
		)

		return strings.HasPrefix(info.Name(), mockPrefix) &&
			strings.HasSuffix(info.Name(), mockSuffix)
	}

	entries, err := fsutils.GatherMatchedEntriesInDir(fs, path, match)
	if err != nil {
		return fmt.Errorf("gather matched entries: %w", err)
	}

	for _, entryPath := range entries {
		if err = fs.Remove(entryPath); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
	}

	return nil
}

func IsCAExists(fs afero.Fs, output string) bool {
	certFilePath := filepath.Join(output, CACertFile)

	_, err := fs.Stat(certFilePath)
	if err == nil {
		return true
	}

	return false
}
