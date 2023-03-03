CLI_NAME=reload

.PHONY: build
build:
	@go build -o ${CLI_NAME} *.go

.PHONY: install
install:
	go mod tidy
	go build -o out *.go
	mv out ${GOBIN}/${CLI_NAME}
