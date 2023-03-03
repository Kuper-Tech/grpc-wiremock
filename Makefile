SHELL=/bin/bash -euo pipefail

.PHONY: check_environment
check_environment:
	@$(shell pwd)/scripts/helpers/check_dependencies.sh

.PHONY: coverage
coverage:
	@go-acc $(shell go list ./... | grep -v -e static) && go tool cover -func=coverage.txt

.PHONY: unit-tests
unit-tests:
	@go test $(shell go list ./... | grep -v -e static) -count 1 -race

.PHONY: install-cli
install-cli:
	make install -C cmd/grpc2http
	make install -C cmd/confgen
	make install -C cmd/watcher
	make install -C cmd/certgen
	make install -C cmd/reload
