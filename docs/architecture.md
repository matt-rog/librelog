# Architecture

LibreLog has three services:

**Web API** - auth, logset CRUD, log queries, and serves the frontend. This is the main service users interact with.

**Ingester** - accepts log data over REST and WebSocket. Separate from the web API so it can be scaled independently for high-throughput use cases.

**Cassandra** - stores everything. Schema is in `cassandra/init.cql`.

## Auth

Account numbers are randomly generated 10-digit numbers. No email, no phone, no PII. Passwords are bcrypt hashed. Account numbers are SHA-256 hashed before storage.

Tokens are also SHA-256 hashed before storage. Session tokens (from login) expire after 30 days via Cassandra TTL. API keys don't expire until revoked.

## Data Model

A logset is a named collection of log entries. Each entry is a JSON object stored as text. No schema enforcement, so each logset can hold whatever shape of data you want.

