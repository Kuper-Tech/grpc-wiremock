package certificates

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/strutils"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

var domains = []string{
	"com", "net", "gov", "mil", "io",
	"ru", "org", "local", "tech", "dev",
	"online", "internal", "team",
}

var commonDomains = []string{
	"localhost", "mock",
	"grpc-wiremock",
}

func (g *certsGen) collectDomains(commonDomains, domains []string) ([]string, error) {
	wiremockConfig, err := g.opener.Open()
	if err != nil {
		log.Printf("certgen: open wiremock config: %s", err)

		if err = handleOpenerErrors(err); err != nil {
			return nil, err
		}
	}

	var targetDomains []string

	for _, service := range wiremockConfig.Services {
		targetDomains = append(targetDomains, service.Name)
	}

	var certDomains []string

	for _, domain := range targetDomains {
		for _, ext := range domains {
			certDomains = append(certDomains, domain, fmt.Sprintf("%s.%s", domain, ext))
		}
	}

	if len(certDomains) == 0 {
		return strutils.UniqueAndSorted(commonDomains...), nil
	}

	return strutils.UniqueAndSorted(append(certDomains, commonDomains...)...), nil
}

func handleOpenerErrors(err error) error {
	if errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, config.EmptyWiremockConfigErr) {
		return nil
	}

	return fmt.Errorf("read wiremock config: %w", err)
}
