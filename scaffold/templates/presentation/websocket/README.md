# WebSocket

WebSocket server implemented using [gorilla/websocket](https://github.com/gorilla/websocket).

Clients connect to `/ws` and exchange JSON messages using an action/payload envelope:

```json
// inbound  (client → server)
{ "action": "get_accounts", "payload": {} }

// outbound (server → client)
{ "action": "get_accounts", "payload": [...] }
```

## File structure

| File | Responsibility |
|---|---|
| `server.go` | Server lifecycle — `Start` / `Stop` |
| `routes.go` | HTTP mux setup, upgrader config, path → handler mapping |
| `handler_account.go` | Account message handler — reads actions, calls use cases, writes responses |
| `common.go` | Shared `message` and `response` envelope types |

## Supported actions

| Action | Description |
|---|---|
| `get_accounts` | Returns a list of accounts |

## Adding a new handler

1. Create `handler_<domain>.go` with a handler struct and a `handle` method.
2. Register the path in `routes.go`:
   ```go
   mux.HandleFunc("/ws/orders", newOrderHandler(ctr, upgrader).handle)
   ```

## Enabling / disabling

Toggle the server without touching code — set `enabled` in `config.yaml`:

```yaml
ws:
  enabled: true
  port: 8003
  wait: 5s
```

## Removing WebSocket support permanently

```
delete  internal/server/server_websocket.go
delete  presentation/websocket/
```

The project will compile and run without it.
