# Service

A Go microservice scaffolded with [Catalyst](https://github.com/kosatnkn/catalyst).

## Project structure

```
cmd/
  server/              Entry point — package main
    main.go            Calls run()
    run.go             Startup orchestration
    config.go          Config defaults and parsing
    migrate.go         Database migration runner
    signal.go          Graceful shutdown signal wait
    fatal.go           Fatal error handler
internal/
  server/              Server protocol registry — package server
    servers.go         ServerLifecycle interface, Register, StartServers
    server_rest.go     REST registration (delete to remove)
    server_graphql.go  GraphQL registration (delete to remove)
    server_grpc.go     gRPC registration (delete to remove)
    server_websocket.go WebSocket registration (delete to remove)
    server_metrics.go  Metrics registration (delete to remove)
  migrations/
    embed.go           Embeds SQL migration files
    sql/               Goose-format migration files
domain/
  entities/            Plain data types with no dependencies
  boundary/            Use-case interfaces (ports)
  usecases/            Business logic
infra/                 Dependency injection container and config
persistence/           Database adapters (implements domain/boundary)
presentation/
  rest/                REST HTTP server (Gin)
  graphql/             GraphQL server (graphql-go)
  grpc/                gRPC server
  websocket/           WebSocket server (gorilla/websocket)
  metrics/             OTel metrics server — exposes /metrics for Prometheus
metadata/              Build-time metadata (populated by set_metadata.sh)
```

## Getting started

```bash
cp config.example.yaml config.yaml   # edit DB credentials and server flags
go mod tidy
make run
```

## Configuration

Each server protocol can be toggled independently in `config.yaml`:

```yaml
rest:
  enabled: true
  port: 8000

graphql:
  enabled: false
  port: 8001

grpc:
  enabled: false
  port: 8002

ws:
  enabled: false
  port: 8003

metrics:
  enabled: false
  port: 9090   # exposes GET /metrics
```

## Removing an unused protocol permanently

Delete the registration file in `internal/server/` and the matching presentation directory:

| Protocol  | Files to delete |
|-----------|-----------------|
| REST      | `internal/server/server_rest.go` + `presentation/rest/` |
| GraphQL   | `internal/server/server_graphql.go` + `presentation/graphql/` |
| gRPC      | `internal/server/server_grpc.go` + `presentation/grpc/` |
| WebSocket | `internal/server/server_websocket.go` + `presentation/websocket/` |
| Metrics   | `internal/server/server_metrics.go` + `presentation/metrics/` |

No other files need to change — the project will compile and run cleanly.

## Common tasks

| Task | Command |
|------|---------|
| Run locally | `make run` |
| Build binary | `make build` |
| Run tests | `make test` |
| Build Docker image | `make docker-build` |
| List upgradable deps | `make dep-upgrade-list` |
| Upgrade all deps | `make dep-upgrade-all` |

## Database migrations

SQL migration files live in `internal/migrations/sql/` and are embedded into the binary at compile time. The service runs all pending migrations automatically on startup using [postgres-migrator](https://github.com/tlmanz/postgres-migrator) (Goose-format).

To add a migration, create a new numbered file:

```
internal/migrations/sql/00002_add_column.sql
```

```sql
-- +goose Up
ALTER TABLE accounts ADD COLUMN email VARCHAR(255);

-- +goose Down
ALTER TABLE accounts DROP COLUMN email;
```
