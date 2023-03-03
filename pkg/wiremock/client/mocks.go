package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/httputils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

func (w *wiremock) UpdateMocks(ctx context.Context, domain string) error {
	port, err := w.findPortByDomain(domain)
	if err != nil {
		return fmt.Errorf("find port by domain '%s': %w", domain, err)
	}

	if err = w.resetMocks(ctx, port); err != nil {
		return fmt.Errorf("reset mocks: %w", err)
	}

	return nil
}

func (w *wiremock) resetMocks(ctx context.Context, port uint) error {
	const path = "__admin/mappings/reset"

	url := fmt.Sprintf("%s:%d/%s", w.host, port, path)

	status, err := httputils.DoPost(ctx, w.client, url, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return httputils.AssertStatus(http.StatusOK, status)
}

func (w *wiremock) findPortByDomain(domain string) (uint, error) {
	wiremockConfig, err := w.configOpener.Open()
	if err != nil {
		return 0, fmt.Errorf("open: %w", err)
	}

	for _, service := range wiremockConfig.Services {
		if service.Name == domain {
			return uint(service.Port), nil
		}
	}

	return 0, config.NoCorrespondingAPIErr
}
