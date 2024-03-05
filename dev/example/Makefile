SHELL=/bin/bash -euo pipefail

MOCKS_PATH="$(shell pwd)/test/wiremock"
CERTS_PATH="$(shell pwd)/certs"
CONTRACTS_PATH="$(shell pwd)/deps"

export GRPC_WIREMOCK_IMAGE=sbermarkettech/grpc-wiremock:dev

compose-up:
	source build/.env && \
		docker compose -f build/docker-compose.yaml up -d

compose-exec:
	@echo -e "Entering to 'wearable-mock' container\nTo develop run ./dev/watch_wiremock.sh\n"
	@source build/.env && \
		docker compose -f build/docker-compose.yaml exec -it --workdir=/source wearable-mock bash

compose-logs:
	source build/.env && \
		docker compose -f build/docker-compose.yaml logs -f --tail=50

compose-up-recreate:
	source build/.env && \
		docker compose -f build/docker-compose.yaml up -d --force-recreate --build

compose-down-all:
	source build/.env && \
		docker compose -f build/docker-compose.yaml down --rmi all

docker-run:
	docker run \
	  -p 9000:9000 \
	  -p 8000:8000 \
	  -v ${MOCKS_PATH}:/home/mock \
	  -v ${CERTS_PATH}:/etc/ssl/mock/share \
	  -v ${CONTRACTS_PATH}:/contracts \
	  ${GRPC_WIREMOCK_IMAGE}
