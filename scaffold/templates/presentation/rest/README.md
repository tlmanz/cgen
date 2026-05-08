# REST

REST server implemented using [Gin](https://github.com/gin-gonic/gin) with
[gin-contrib/cors](https://github.com/gin-contrib/cors) middleware.

## File structure

| File | Responsibility |
|---|---|
| `server.go` | Server lifecycle — `Start` / `Stop` |
| `routes.go` | Router setup, middleware registration, path → controller mapping |
| `controller_info.go` | Info, health (`/healthz`), and readiness (`/readyz`) endpoints |
| `controller_account.go` | Account endpoints — `GET /accounts`, `POST /accounts` |
| `middleware.go` | Logger middleware |
| `config.go` | Gin debug/release mode configuration |
| `common_requests.go` | Request parsing helpers (paging, filters, request structs) |
| `common_responses.go` | Standardised response envelope (`data`, `error`) |
| `helpers.go` | Shared handler utilities |

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Service info |
| `GET` | `/healthz` | Kubernetes liveness probe |
| `GET` | `/readyz` | Kubernetes readiness probe |
| `GET` | `/accounts` | List accounts (supports `?paging={"page":1,"size":10}&filters={}`) |
| `POST` | `/accounts` | Create an account |

## Adding a new controller

1. Create `controller_<domain>.go` with a controller struct and handler methods.
2. Add a `registerXxxHandlers` function and call it from `routes.go`.

## Enabling / disabling

Toggle the server without touching code — set `enabled` in `config.yaml`:

```yaml
rest:
  enabled: true
  port: 8000
  wait: 5s
  release: false  # set to true in production
```

## Removing REST support permanently

```
delete  internal/server/server_rest.go
delete  presentation/rest/
```

The project will compile and run without it.
