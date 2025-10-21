# Configuration

All config is via environment variables.

## Web API

| Variable | Description | Default |
|---|---|---|
| `CASSANDRA_CLUSTER` | Cassandra host(s), space-separated | `librelog-cassandra` |
| `PUBLIC_REGISTRATION` | Allow new signups | `false` |

## Ingester

| Variable | Description | Default |
|---|---|---|
| `CASSANDRA_CLUSTER` | Cassandra host(s), space-separated | `librelog-cassandra` |

## Private Nodes

By default, registration is closed. To add accounts, temporarily set `PUBLIC_REGISTRATION=true`, create your accounts, then set it back to `false`. Or sign up once and create API keys for your devices.
