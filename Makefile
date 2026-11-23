GOBIN              ?= $(shell go env GOPATH)/bin
MOCKERY            ?= $(GOBIN)/mockery
SWAG               ?= $(GOBIN)/swag
COVERAGE_EXCLUDE   ?= mocks|main.go|test|docs/docs.go
COVERAGE_THRESHOLD ?= 80

.PHONY: tools mock swag test test-unit test-integration run redis-up redis-down docker-up docker-down

tools:
	go install github.com/vektra/mockery/v3@v3.9.0
	go install github.com/swaggo/swag/cmd/swag@v1.16.6

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

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down
