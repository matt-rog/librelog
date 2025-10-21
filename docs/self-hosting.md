# Self-Hosting

## Requirements

Docker and Docker Compose.

## Quick Start

```
git clone https://github.com/matt-rog/librelog
cd librelog
docker compose up -d
```

This starts Cassandra, runs the schema migration, then brings up the web API and ingester.

## Production

Put a reverse proxy (Caddy, nginx) in front for TLS. Cassandra data persists in a Docker volume (`cassandra-data`).

## Backups

Snapshot the Cassandra Docker volume. To restore, replace the volume.
