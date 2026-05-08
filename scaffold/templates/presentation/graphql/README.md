# GraphQL

GraphQL server implemented using [graphql-go](https://github.com/graphql-go/graphql) with the
[graphql-go/handler](https://github.com/graphql-go/handler) HTTP adapter.
The GraphQL Playground UI is available at `/graphql` in a browser.

## File structure

| File | Responsibility |
|---|---|
| `server.go` | Server lifecycle — `Start` / `Stop` |
| `schema.go` | Assembles the full schema from query and mutation types |
| `types.go` | GraphQL type definitions (e.g. `Account`) |
| `resolver_account.go` | Account queries (`accounts`) and mutations (`createAccount`) |

## Adding a new domain

1. Define new GraphQL types in `types.go`.
2. Create `resolver_<domain>.go` with `newXxxQueryType` / `newXxxMutationType` functions.
3. Wire them into `schema.go`.

## Sample queries

```graphql
# Get accounts
query {
  accounts {
    id
    owner
    currency
    balance
  }
}

# Create an account
mutation {
  createAccount(owner: "Alice", currency: "USD") {
    id
    owner
    currency
    balance
  }
}
```

## Enabling / disabling

Toggle the server without touching code — set `enabled` in `config.yaml`:

```yaml
graphql:
  enabled: true
  port: 8001
  wait: 5s
```

## Removing GraphQL support permanently

```
delete  internal/server/server_graphql.go
delete  presentation/graphql/
```

The project will compile and run without it.
