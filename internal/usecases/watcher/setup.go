package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/certgen"
	"github.com/SberMarket-Tech/grpc-wiremock/internal/usecases/confgen"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/environment"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/runner"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/watcher"
	wmclient "github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/client"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/configopener"
)

const (
	MocksWatcher   = "mocks"
	DomainsWatcher = "domains"
)

const (
	filesDirRule    = `.*/__files$`
	mappingsDirRule = `.*/mappings$`

	mockFiles = `^.*\/(mappings|__files)\/.+[.]json$`
)

type WatchRequest struct {
	Name string
	Path string
}

var (
	osfs = afero.NewOsFs()

	opener = configopener.New(osfs, environment.SupervisordConfigsDirPath)

	wiremockClient = wmclient.NewDefaultClient(osfs, opener)
)

var knownWatchers = map[string]watcher.WatcherDesc{
	MocksWatcher: {
		Do: mocksWatcherAction,

		Name: MocksWatcher,

		Recursive: true,

		Behave: watcher.BehaviourDesc{
			Event: watcher.NewEventTypes().
				WithCreate().WithRemove().WithWrite(),

			Throttle: watcher.ThrottlingRules{
				DelayAfterEvent: time.Second * 3,
				Interval:        time.Millisecond * 500,
			},

			Retry: watcher.RetryRules{Attempts: 10},

			Entry: watcher.NewEntryRules().
				WithNameRule(mockFiles).
				WithNameRule(filesDirRule).
				WithNameRule(mappingsDirRule),
		},
	},

	//Turn off domains' watcher.

	//DomainsWatcher: {
	//	Do: domainsWatcherAction,
	//
	//	Name: DomainsWatcher,
	//
	//	Recursive: false,
	//
	//	Behave: watcher.BehaviourDesc{
	//		Event: watcher.NewEventTypes().
	//			WithCreate().
	//			WithRemove().
	//			WithWrite().
	//			WithRename(),
	//
	//		Retry: watcher.RetryRules{Attempts: 10},
	//
	//		Entry: watcher.NewEntryRules().WithNameRule(`.*`),
	//	},
	//},
}

func mocksWatcherAction(ctx context.Context, _ string, path string) error {
	// only for '.json' files in '/home/mock/*' directory.

	domain, err := getDomainByPath(path)
	if err != nil {
		return nil
	}

	log.Printf("update mocks for domain '%s'\n", domain)

	if err = wiremockClient.UpdateMocks(ctx, domain); err != nil {
		return fmt.Errorf("update mocks for domain '%s': %w", domain, err)
	}

	return nil
}

func domainsWatcherAction(ctx context.Context, root string, path string) error {
	/*
		Wiremock directories:
		/home/mock/
			domain1
				__files
				mappings
			domain2
				__files
				mappings

		Supervisord directories:
		/etc/supervisord/
			mocks/
				domain1-mock.conf
				domain2-mock.conf
			supervisord.conf

		NGINX directories:
		/etc/nginx/http.d/
			domain1-mock.conf
			domain2-mock.conf
			default.conf
	*/

	domain, err := getDomainByFolder(root, path)
	if err != nil {
		log.Printf("get domain by folder: %s\n", err.Error())
		return nil
	}

	commandRunner := runner.New(exec.New())

	configsGenerator := confgen.NewConfGenWithDefaultFs(
		environment.DefaultWiremockConfigPath,
		environment.SupervisordConfigsDirPath,
		commandRunner, os.Stdout)

	opts := confgen.Options{
		confgen.WithNGINX(environment.NginxConfigsPath),
		confgen.WithSupervisord(environment.SupervisordConfigsPath),
	}

	if err = configsGenerator.Generate(ctx, opts...); err != nil {
		return fmt.Errorf("generate configs: %w", err)
	}

	log.Println("update supervisord & nginx config", domain)

	certGenerator := certgen.NewCertsGenWithDefaultFs(
		environment.DefaultCertificatesPath,
		environment.SupervisordConfigsDirPath,
		os.Stdout,
	)

	log.Println("update certificates", domain)

	if err := certGenerator.Generate(ctx); err != nil {
		return fmt.Errorf("generate certificates: %w", err)
	}

	return nil
}

var domainRE = regexp.MustCompile(`(.*/)?(.*)/(mappings|__files)`)

func getDomainByPath(path string) (string, error) {
	domainSubmatches := domainRE.FindAllStringSubmatch(path, 1)

	if len(domainSubmatches) == 0 || len(domainSubmatches[0]) < 3 {
		return "", fmt.Errorf("get domain by path '%s'", path)
	}

	domain := domainSubmatches[0][2]

	return domain, nil
}

func getDomainByFolder(rootPath, eventPath string) (string, error) {
	relativePath, err := filepath.Rel(rootPath, eventPath)
	if err != nil {
		return "", fmt.Errorf("relative path: %w", err)
	}

	if !isFirstLevelFolder(relativePath) {
		return "", fmt.Errorf("path is not correspond to first level folder: %s", relativePath)
	}

	return filepath.Base(relativePath), nil
}

func isFirstLevelFolder(path string) bool {
	if len(filepath.Ext(path)) > 0 {
		return false
	}

	return len(strings.Split(path, string(os.PathSeparator))) == 1
}
