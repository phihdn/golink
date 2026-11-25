GOBIN              ?= $(shell go env GOPATH)/bin
MOCKERY            ?= $(GOBIN)/mockery
SWAG               ?= $(GOBIN)/swag
COVERAGE_EXCLUDE   ?= mocks|main.go|test|docs/docs.go
COVERAGE_THRESHOLD ?= 80
IMAGE              ?= golink:local
DOCKER_PORT        ?= 8080
# Same formatters, order and versions as the editor (conform.nvim): gofumpt, golines, goimports-reviser.
GOFUMPT            ?= go run mvdan.cc/gofumpt@v0.12.0
GOLINES            ?= go run github.com/golangci/golines@v0.15.0
GOIMPORTS_REVISER  ?= go run github.com/incu6us/goimports-reviser/v3@v3.13.2
# Hand-written Go files. Generated code (swag docs, mockery mocks) is left as its generator wrote it.
GO_FILES            = $(shell find . -name '*.go' -not -path './docs/*' -not -path '*/mocks/*')

.PHONY: tools hooks fmt mock swag test test-unit test-integration run redis-up redis-down docker-build docker-run docker-up docker-down

tools:
	go install github.com/vektra/mockery/v3@v3.9.0
	go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Points git at .githooks, so the pre-commit hook formats code before every commit.
hooks:
	git config core.hooksPath .githooks

fmt:
	$(GOFUMPT) -w $(GO_FILES)
	$(GOLINES) -w $(GO_FILES)
	@$(GOIMPORTS_REVISER) -format $(GO_FILES) 2>/dev/null

mock:
	$(MOCKERY)

swag:
	$(SWAG) init -g cmd/api/main.go -o docs --parseInternal

# Runs all tests with coverage and fails when total coverage is below COVERAGE_THRESHOLD.
test: mock
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

test-unit:
	go test $$(go list ./... | grep -v /test/)

test-integration:
	go test ./internal/test/...

run:
	go run ./cmd/api

redis-up:
	docker compose up -d redis

redis-down:
	docker compose stop redis

# Builds the multi-stage Dockerfile into IMAGE.
docker-build:
	docker build -t $(IMAGE) .

# Runs IMAGE on DOCKER_PORT against the Redis from redis-up on the host.
docker-run: docker-build redis-up
	docker run --rm --name golink -p $(DOCKER_PORT):8080 \
		--add-host=host.docker.internal:host-gateway \
		-e REDIS_ADDR=host.docker.internal:6379 \
		$(IMAGE)

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down
