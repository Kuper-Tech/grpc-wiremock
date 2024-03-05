#!/bin/bash

set -euo pipefail

echo docker build . \
    --build-arg GOLANG_IMAGE_REPO=registry.hub.docker.com \
    --build-arg GOLANG_IMAGE_NAME=library/golang \
    --build-arg GOLANG_IMAGE_TAG=1.21.7 \
    --build-arg WIREMOCK_IMAGE_REPO=docker.io \
    --build-arg WIREMOCK_IMAGE_NAME=wiremock/wiremock \
    --build-arg WIREMOCK_IMAGE_TAG=2.32.0-alpine \
    --tag=sbermarkettech/grpc-wiremock:dev

docker build . \
    --build-arg GOLANG_IMAGE_REPO=registry.hub.docker.com \
    --build-arg GOLANG_IMAGE_NAME=library/golang \
    --build-arg GOLANG_IMAGE_TAG=1.21.7 \
    --build-arg WIREMOCK_IMAGE_REPO=docker.io \
    --build-arg WIREMOCK_IMAGE_NAME=wiremock/wiremock \
    --build-arg WIREMOCK_IMAGE_TAG=2.32.0-alpine \
    --tag=sbermarkettech/grpc-wiremock:dev
