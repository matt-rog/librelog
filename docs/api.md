# API Reference

All endpoints return JSON. Authenticate with `Authorization: Bearer <token>`.

## Quick Start

```
# sign up (if registration is open)
ACCT=$(curl -s -X POST localhost:8080/api/signup \
  -d '{"password":"mypass"}' | jq -r .account_number)

# login
TOKEN=$(curl -s -X POST localhost:8080/api/login \
  -d "{\"account_number\":\"$ACCT\",\"password\":\"mypass\"}" | jq -r .token)

# create a logset
ID=$(curl -s -X POST localhost:8080/api/logsets \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"weight"}' | jq -r .log_id)

# ingest a log entry
curl -X POST localhost:9000/ingest \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"log_set\":\"$ID\",\"data\":{\"kg\":81.2}}"

# read it back
curl "localhost:8080/api/logsets/$ID/logs" \
  -H "Authorization: Bearer $TOKEN"
```

## Auth

### POST /api/signup

```
curl -X POST localhost:8080/api/signup \
  -d '{"password": "s3cret", "name": "optional"}'
```

```
{"account_number": "3847291056"}
```

Save the account number. It's only shown once.

### POST /api/login

```
curl -X POST localhost:8080/api/login \
  -d '{"account_number": "3847291056", "password": "s3cret"}'
```

```
{"token": "a57a8f7f..."}
```

### POST /api/logout

```
curl -X POST localhost:8080/api/logout \
  -H "Authorization: Bearer $TOKEN"
```

## Logsets

### GET /api/logsets

```
curl localhost:8080/api/logsets \
  -H "Authorization: Bearer $TOKEN"
```

```
[{"log_id": "abc-123", "name": "running", "description": ""}]
```

### POST /api/logsets

```
curl -X POST localhost:8080/api/logsets \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "running", "description": "daily runs"}'
```

### PUT /api/logsets/:id

```
curl -X PUT localhost:8080/api/logsets/abc-123 \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "running-2025"}'
```

### DELETE /api/logsets/:id

```
curl -X DELETE localhost:8080/api/logsets/abc-123 \
  -H "Authorization: Bearer $TOKEN"
```

## Logs

### GET /api/logsets/:id/logs

Params: `limit` (1-1000, default 100), `before` / `after` (RFC3339 timestamps for pagination).

```
curl "localhost:8080/api/logsets/abc-123/logs?limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

```
[{"recv_time": "2025-10-13T20:00:00Z", "data": "{\"miles\": 3.2}"}]
```

### GET /api/logsets/:id/export

Params: `format` - `json` (default) or `csv`. Exports all entries.

```
curl "localhost:8080/api/logsets/abc-123/export?format=csv" \
  -H "Authorization: Bearer $TOKEN" -o running.csv
```

## Ingesting Data

### POST /ingest

```
curl -X POST localhost:9000/ingest \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"log_set": "abc-123", "data": {"miles": 3.2, "time_min": 28}}'
```

```
{"status": "ok"}
```

### WebSocket /ingest

Connect with token as query param. Send JSON messages, get `{"status":"ok"}` back for each.

```
import asyncio, websockets, json

async def main():
    uri = "ws://localhost:9000/ingest?token=" + TOKEN
    async with websockets.connect(uri) as ws:
        await ws.send(json.dumps({
            "log_set": "abc-123",
            "data": {"miles": 3.2}
        }))
        print(await ws.recv())

asyncio.run(main())
```

## API Keys

Long-lived tokens for scripts, cron jobs, and plugins. No expiry until revoked.

### POST /api/tokens

```
curl -X POST localhost:8080/api/tokens \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "vscode-laptop"}'
```

```
{"token": "9d3fd1ec...", "name": "vscode-laptop", "prefix": "9d3fd1ec"}
```

The full token is only shown once.

### GET /api/tokens

```
curl localhost:8080/api/tokens \
  -H "Authorization: Bearer $TOKEN"
```

### DELETE /api/tokens/:hash

```
curl -X DELETE localhost:8080/api/tokens/<token_hash> \
  -H "Authorization: Bearer $TOKEN"
```
