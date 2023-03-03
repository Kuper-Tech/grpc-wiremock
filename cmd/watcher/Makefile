CLI_NAME=watcher

.PHONY: build
build:
	@go build -o ${CLI_NAME} *.go

.PHONY: install
install:
	go mod tidy
	go build -o out *.go
	mv out ${GOBIN}/${CLI_NAME}
