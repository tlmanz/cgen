# gRPC

gRPC server implemented using [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc).

The sample uses a **JSON codec** (registered in `codec.go`) in place of the default protobuf codec
so the server runs without any code generation step. The proto definition is kept in
`proto/account.proto` as a reference for when you are ready to switch to real protobuf.

## File structure

| File | Responsibility |
|---|---|
| `server.go` | Server lifecycle — `Start` / `Stop` |
| `routes.go` | Service descriptor, `AccountServiceServer` interface, handler registration |
| `service_account.go` | Account service implementation wired to use cases |
| `common.go` | Request / response types |
| `codec.go` | JSON codec — replaces the proto codec for this sample |
| `proto/account.proto` | Reference proto definition |

## Adding a new service

1. Add request / response types to `common.go`.
2. Create `service_<domain>.go` implementing the new service interface.
3. Add a `grpc.ServiceDesc` and registration function to `routes.go`.
4. Call `RegisterXxxServer(s.srv, ...)` in `server.go`.

## Migrating to real protobuf

When you are ready to use proper protobuf instead of the JSON codec:

1. Delete `codec.go`.
2. Install the code generation toolchain:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
3. Generate Go code from `proto/account.proto`:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/account.proto
   ```
4. Replace the hand-written types in `common.go` and descriptor in `routes.go`
   with the generated `*pb` types.

## Enabling / disabling

Toggle the server without touching code — set `enabled` in `config.yaml`:

```yaml
grpc:
  enabled: true
  port: 8002
```

## Removing gRPC support permanently

```
delete  internal/server/server_grpc.go
delete  presentation/grpc/
```

The project will compile and run without it.
