OUTPUT=output

.PHONY: clean
clean:
	@rm -rf "${OUTPUT}" && mkdir "${OUTPUT}"

.PHONY: up
up:
	@WIREMOCK_API_COUNT=$(WIREMOCK_API_COUNT) && docker compose up --force-recreate --build

.PHONY: down
down:
	docker compose down

.PHONY: benchmark
benchmark:
	bash benchmark.sh
