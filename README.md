# golink

URL shortener service, exercise project for EBVN06. Built with Go, Echo v5, Redis and a clean architecture layout.

## Project layout

| Path | Purpose |
| --- | --- |
| `cmd/api` | Entry point. Loads config and starts the server. |
| `internal/config` | Environment based config (viper). |
| `internal/api` | Echo setup and route registration. |
| `internal/handler` | HTTP handlers. |
| `internal/service` | Business logic. |
| `internal/repository` | Redis access. |
| `internal/model` | Request and response types. |
| `internal/test/integration_test` | End to end tests through the full router. |
| `docs` | Generated Swagger spec. |

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port. |
| `SERVICE_NAME` | `golink` | Name returned by the health check. |
| `INSTANCE_ID` | random UUID | Instance id returned by the health check. |
| `REDIS_ADDR` | `localhost:6379` | Redis address. |
| `REDIS_PASSWORD` | empty | Redis password. |
| `REDIS_DB` | `0` | Redis database number. |

## Getting started

```bash
make tools      # install mockery and swag
make redis-up   # start Redis in Docker
make run        # start the API on :8080
```

Run the API and Redis together in Docker with `make docker-up`, and stop them with `make docker-down`.

## Make targets

| Target | Description |
| --- | --- |
| `tools` | Install mockery and swag. |
| `mock` | Regenerate mocks from `.mockery.yml`. |
| `swag` | Regenerate the Swagger spec in `docs/`. |
| `test` | Run all tests with coverage. Writes `coverage.out` and `coverage.html` and fails below `COVERAGE_THRESHOLD` (default 80). |
| `test-unit` | Run unit tests only. |
| `test-integration` | Run integration tests only. |
| `run` | Start the API locally. |
| `redis-up` / `redis-down` | Start or stop Redis in Docker. |
| `docker-up` / `docker-down` | Start or stop the API and Redis in Docker. |

## API

Swagger UI is served at `http://localhost:8080/swagger/index.html`.

### GET /health-check

Returns `200` with the service name and instance id.

```json
{"message": "OK", "service_name": "golink", "instance_id": "1f0c..."}
```

### GET /v1/gen-pass

Returns a random 10-character alphanumeric password as plain text.
