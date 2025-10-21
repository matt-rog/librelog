![LibreLog](docs/librelog.png)

Most apps that hold your personal data don't let you really own it. LibreLog is an open-source, self-hostable log store. Send data in, pull it out, do whatever you want with it.

## Install

```
git clone https://github.com/matt-rog/librelog
cd librelog
docker compose up -d
```

See [self-hosting](docs/self-hosting.md) and [configuration](docs/configuration.md) for production setup.

## Who is this for?

People who want a simple, centralized place to log data without building custom infrastructure every time. LibreLog is kept minimal on purpose to  meet consumer device constraints, encouraging self-hosting. This also helps make plugin/connector development easier. 

If you need serious analytics or high-throughput, you'd probably want to bring your own logging solution.

## Docs

- [Architecture](docs/architecture.md)
- [Self-Hosting](docs/self-hosting.md)
- [Configuration](docs/configuration.md)
- [API Reference](docs/api.md)
