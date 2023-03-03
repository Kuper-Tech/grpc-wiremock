#!/bin/bash

set -euo pipefail

echo "===== Dependencies check is started..."

if command -v protoc > /dev/null; then
  echo "- proto compiler [OK]"
else
  echo "- proto compiler [FAIL]"
  echo "To install proto compiler: https://grpc.io/docs/protoc-installation"
fi

GO="protoc-gen-go"
if command -v ${GO} > /dev/null; then
  echo "- ${GO} plugin [OK]"
else
  echo "- ${GO} plugin [FAIL]"
  echo "To install ${GO} plugin:"
  echo "- go install google.golang.org/protobuf/cmd/${GO}@latest"
fi

GO_GRPC="protoc-gen-go-grpc"
if command -v ${GO_GRPC} > /dev/null; then
  echo "- ${GO_GRPC} plugin [OK]"
else
  echo "- ${GO_GRPC} plugin [FAIL]"
  echo "To install ${GO_GRPC} plugin:"
  echo "- go install google.golang.org/grpc/cmd/${GO_GRPC}@latest"
fi

OPENAPI="protoc-gen-openapi"
if command -v ${OPENAPI} > /dev/null; then
  echo "- ${OPENAPI} plugin [OK]"
else
  echo "- ${OPENAPI} plugin [FAIL]"
  echo "To install ${OPENAPI} plugin:"
  echo "- https://github.com/solo-io/${OPENAPI}"
fi

echo "===== Dependencies check is done."