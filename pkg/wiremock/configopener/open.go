package configopener

import (
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	supervisorconf "github.com/ochinchina/supervisord/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/wiremock/config"
)

type opener struct {
	fs afero.Fs

	path string
}

func New(fs afero.Fs, path string) *opener {
	return &opener{fs: fs, path: path}
}

func (o *opener) Open() (config.Wiremock, error) {
	const configName = "supervisord.conf"

	var wiremock config.Wiremock

	path := filepath.Join(o.path, configName)
	supervisordConf := supervisorconf.NewConfig(path)

	logrus.SetOutput(io.Discard)

	_, err := supervisordConf.Load()
	if err != nil {
		return wiremock, err
	}

	for _, p := range supervisordConf.GetPrograms() {
		envs := convertEnvs(p.GetEnv("environment"))

		name, exists := envs["NAME"]
		if !exists {
			continue
		}

		root, exists := envs["ROOT"]
		if !exists {
			continue
		}

		port, exists := envs["PORT"]
		if !exists {
			continue
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			continue
		}

		wiremock.Services = append(wiremock.Services, config.Service{
			Name: name, RootDir: root, Port: portInt})
	}

	return wiremock, nil
}

func convertEnvs(envs []string) map[string]string {
	const separator = "="
	converted := map[string]string{}

	for _, env := range envs {
		parts := strings.Split(env, separator)
		if len(parts) != 2 {
			log.Println("parse env error:", env)
			continue
		}

		converted[parts[0]] = parts[1]
	}

	return converted
}
