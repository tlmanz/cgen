# Metrics

OpenTelemetry metrics exported in Prometheus format via a dedicated HTTP server.

The metrics server exposes `GET :<port>/metrics` which Prometheus scrapes.
Go runtime metrics (goroutines, GC, heap) are included automatically by the Prometheus exporter.

## File structure

| File | Responsibility |
|---|---|
| `server.go` | Server lifecycle — sets up OTel provider + Prometheus exporter, serves `/metrics` |
| `middleware.go` | Gin middleware recording HTTP request count and latency per route |
| `business.go` | Sample business counter (`accounts_created_total`) — shows the pattern for domain metrics |

## Enabling / disabling

Toggle the server without touching code — set `enabled` in `config.yaml`:

```yaml
metrics:
  enabled: true
  port: 9090   # Prometheus scrapes GET :9090/metrics
```

## Wiring the HTTP middleware into REST

The HTTP request count and latency middleware is opt-in. Add it to your router in
`presentation/rest/routes.go`:

```go
import "yourmodule/presentation/metrics"

func newRouter(cfg infra.RESTConfig, ctr *infra.Container) *gin.Engine {
    router := gin.New()
    router.Use(metrics.NewHTTPMiddleware())  // ← add this line
    // ...
}
```

If the metrics server is disabled (`enabled: false`), OTel uses its no-op provider and
the middleware records nothing — no performance cost, no errors.

## Adding a business metric

Follow the pattern in `business.go`. Use `sync.Once` for lazy initialisation so the
instrument resolves against the real OTel provider (set up during server start):

```go
var (
    ordersPlacedOnce    sync.Once
    ordersPlacedCounter metric.Int64Counter
)

func RecordOrderPlaced(ctx context.Context) {
    ordersPlacedOnce.Do(func() {
        meter := otel.Meter("business")
        ordersPlacedCounter, _ = meter.Int64Counter("orders_placed_total",
            metric.WithDescription("Total number of orders placed"))
    })
    ordersPlacedCounter.Add(ctx, 1)
}
```

Call `metrics.RecordOrderPlaced(ctx)` from the relevant use case.

## Migrating to real protobuf / OTLP

To push metrics to an OTel Collector instead of exposing a Prometheus endpoint:

1. Replace `go.opentelemetry.io/otel/exporters/prometheus` with the OTLP exporter:
   ```bash
   go get go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc
   ```
2. Update `server.go` to use `otlpmetricgrpc.New(ctx)` instead of `promexporter.New()`.
3. Remove `github.com/prometheus/client_golang` from `go.mod` if no longer needed.

## Removing metrics support permanently

```
delete  internal/server/server_metrics.go
delete  presentation/metrics/
```

The project will compile and run without it.
