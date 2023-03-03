PROXY_BINARY_NAME=grpc-to-http-proxy

run:
	@go mod tidy
	@go run cmd/*.go

install:
	@go mod tidy
	@go build -o ${PROXY_BINARY_NAME} cmd/*.go
	@mv ${PROXY_BINARY_NAME} ${GOBIN}/${PROXY_BINARY_NAME}
