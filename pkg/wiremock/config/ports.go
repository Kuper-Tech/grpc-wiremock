package config

import (
	"sort"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/sliceutils"
)

type Ports []int

func GatherPorts(wiremock Wiremock) Ports {
	var ports Ports

	for _, service := range wiremock.Services {
		ports = append(ports, service.Port)
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i] < ports[j]
	})

	return ports
}

func (p Ports) Allocate() int {
	const defaultPort = 8000

	if len(p) == 0 {
		return defaultPort
	}

	var ports Ports
	for _, port := range p {
		ports = append(ports, port)
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i] > ports[j]
	})

	return sliceutils.FirstOf(ports) + 1
}
