MOCKERY ?= $(shell go env GOPATH)/bin/mockery

.PHONY: mock test run

mock:
	$(MOCKERY)

test: mock
	go test ./...

run:
	go run ./cmd/api
