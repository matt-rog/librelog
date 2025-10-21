![LibreLog](docs/librelog.png)

Most software that tracks your data doesn't really let you own it. You can't easily export it, query it programmatically, or move it somewhere else. LibreLog is an open-source, self-hostable log store that aims to accept data from anywhere and let you use it however you want.

## Getting Started

When you sign up, you get a randomly generated account number instead of using an email or phone number. **Store it somewhere safe** because it's the only time you'll see it, and your code + your password is the only way to log in.

Once logged in, create a logset (a named collection of log entries) and start pushing data to it. To send data from external apps or scripts, create an API key from the keys panel. API keys won't expire until you revoke them.

See the [API reference](docs/api.md) for details on ingesting and querying data.

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
