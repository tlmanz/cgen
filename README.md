# cgen

A CLI tool that scaffolds production-ready Go microservices following Clean Architecture.

> **Built on [kosatnkn/catalyst](https://github.com/kosatnkn/catalyst)**
> — a Clean Architecture microservice template written in Go by [@kosatnkn](https://github.com/kosatnkn).
> `cgen` packages that template into an installable CLI so you can scaffold new services in one command.
> All architectural decisions, layer structure, and design patterns originate from the original Catalyst project.
> Full credit to the author.

## Installation

```bash
go install github.com/tlmanz/cgen@latest
```

## Usage

```bash
cgen new --module github.com/yourorg/yourservice
```

| Flag | Description | Default |
|------|-------------|---------|
| `--module` | Go module path (required) | — |
| `--dir` | Output directory | `.` (current directory) |
| `--yes` | Skip confirmation prompt | false |

The project directory name is inferred from the last segment of the module path.

## What gets generated

```
yourservice/
  cmd/
    server/              Entry point — package main
      main.go            Calls run()
      run.go             Startup orchestration
      config.go          Config defaults and parsing
      migrate.go         Database migration runner
      signal.go          Graceful shutdown signal wait
      fatal.go           Fatal error handler
  internal/
    server/              Protocol registry — package server
      servers.go         ServerLifecycle interface, Register, StartServers
      server_rest.go     REST registration
      server_graphql.go  GraphQL registration
      server_grpc.go     gRPC registration
      server_websocket.go WebSocket registration
      server_metrics.go  Metrics registration
    migrations/
      embed.go           Embeds SQL files into the binary
      sql/               Goose-format migration files
  domain/
    entities/            Plain data types
    boundary/            Use-case interfaces (ports)
    usecases/            Business logic
  infra/                 DI container and config
  persistence/           Database adapters
  presentation/
    rest/                REST server (Gin)
    graphql/             GraphQL server (graphql-go)
    grpc/                gRPC server (JSON codec, no protoc required)
    websocket/           WebSocket server (gorilla/websocket)
    metrics/             OTel metrics — exposes /metrics for Prometheus
  metadata/              Build-time metadata
  Dockerfile
  Makefile
  config.example.yaml
```

## Getting started after scaffold

```bash
cd yourservice
go mod tidy
cp config.example.yaml config.yaml   # edit DB credentials and ports
make run
```

## Architecture

Each generated service follows Clean Architecture with strict layer boundaries:

```
Presentation  →  Domain (boundary interfaces)  ←  Persistence
                       ↓
                   Use Cases
                       ↓
                    Entities
```

- **Domain** has zero external dependencies — only plain Go.
- **Persistence** implements domain boundary interfaces (ports).
- **Presentation** calls use cases through the container; it never touches persistence directly.
- **infra** wires everything together at startup.

## Protocol servers

Each protocol is independently toggleable. No code changes required:

```yaml
# config.yaml
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
  port: 9090   # exposes GET /metrics for Prometheus
```

### Removing a protocol permanently

Delete two things — nothing else needs to change:

| Protocol  | Files to delete |
|-----------|-----------------|
| REST      | `internal/server/server_rest.go` + `presentation/rest/` |
| GraphQL   | `internal/server/server_graphql.go` + `presentation/graphql/` |
| gRPC      | `internal/server/server_grpc.go` + `presentation/grpc/` |
| WebSocket | `internal/server/server_websocket.go` + `presentation/websocket/` |
| Metrics   | `internal/server/server_metrics.go` + `presentation/metrics/` |

This works because each protocol registers itself via `init()` in its own file. Deleting the file removes it from the registry without touching anything else.

## Included sample implementations

Each protocol ships with a working `accounts` example wired end-to-end through the clean architecture layers:

| Protocol | Sample |
|----------|--------|
| REST | `GET /accounts`, `POST /accounts` |
| GraphQL | `accounts` query, `createAccount` mutation |
| gRPC | `AccountService` with `GetAccounts` and `CreateAccount` methods |
| WebSocket | `get_accounts` action over a JSON message envelope |
| Metrics | HTTP request count + latency middleware, `accounts_created_total` business counter |

### gRPC without protoc

The gRPC sample uses a JSON codec registered at startup (`presentation/grpc/codec.go`), so it runs without any code generation step. A reference `proto/account.proto` is included for when you are ready to migrate to real protobuf.

## Database migrations

SQL files live in `internal/migrations/sql/` and are embedded into the binary at compile time. Migrations run automatically on every startup using [postgres-migrator](https://github.com/tlmanz/postgres-migrator) (Goose-format).

```
internal/migrations/sql/00001_init.sql
internal/migrations/sql/00002_your_change.sql
```

```sql
-- +goose Up
ALTER TABLE accounts ADD COLUMN email VARCHAR(255);

-- +goose Down
ALTER TABLE accounts DROP COLUMN email;
```

## Credits

This tool is a CLI wrapper around **[kosatnkn/catalyst](https://github.com/kosatnkn/catalyst)** — a Clean Architecture microservice template for Go, created and maintained by [@kosatnkn](https://github.com/kosatnkn).

All of the following originate from the Catalyst project:

- The Clean Architecture layer structure (`domain`, `infra`, `persistence`, `presentation`)
- The DI container and config patterns in `infra/`
- The REST server implementation using [Gin](https://github.com/gin-gonic/gin)
- The use of [catalyst-pkgs](https://github.com/kosatnkn/catalyst-pkgs) for logging and configuration
- The `metadata` package and `set_metadata.sh` build tooling
- The `Makefile`, `Dockerfile`, and project conventions

`cgen` adds:

- A Go-installable CLI (`go install`) replacing the original shell-script scaffolding
- GraphQL, gRPC, WebSocket, and Metrics server implementations
- The self-registration pattern (`init()` + `internal/server/`) for protocol plug-and-unplug
- Embedded database migrations via [postgres-migrator](https://github.com/tlmanz/postgres-migrator)
- The `cmd/server/` + `internal/` project layout

If you find Catalyst useful, please star the original repository:
**[github.com/kosatnkn/catalyst](https://github.com/kosatnkn/catalyst)**

## Common make targets

| Command | Description |
|---------|-------------|
| `make run` | Run locally (also sets build metadata) |
| `make build` | Build binary to `./main` |
| `make test` | Run tests with coverage |
| `make docker-build` | Build Docker image |
| `make dep-upgrade-list` | List upgradable dependencies |
| `make dep-upgrade-all` | Upgrade all dependencies |
