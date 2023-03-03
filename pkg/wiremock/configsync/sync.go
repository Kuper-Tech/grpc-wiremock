package configsync

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

func SyncWiremockConfig(fs afero.Fs, wiremockConfig config.Wiremock, wiremockPath string) (config.Wiremock, error) {
	domainDirectories, err := fsutils.GatherDirs(fs, wiremockPath)
	if err != nil {
		return config.Wiremock{}, fmt.Errorf("gather domain directories: %w", err)
	}

	nextAvailablePort := config.GatherPorts(wiremockConfig).Allocate()

	domainToService := map[string]config.Service{}
	targetDomainToService := map[string]config.Service{}

	for _, service := range wiremockConfig.Services {
		domainToService[service.Name] = service
	}

	for _, domainDir := range domainDirectories {
		domain := filepath.Base(domainDir.Name())

		service, exists := domainToService[domain]
		if exists {
			targetDomainToService[domain] = service
		} else {
			targetDomainToService[domain] = config.NewService(wiremockPath, domain, nextAvailablePort)
			nextAvailablePort++
		}
	}

	var targetWiremockConfig config.Wiremock
	for _, service := range targetDomainToService {
		targetWiremockConfig.Services = append(targetWiremockConfig.Services, service)
	}

	return targetWiremockConfig, nil
}
