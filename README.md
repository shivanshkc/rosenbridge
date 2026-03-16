# Rosenbridge

Rosenbridge is a minimal, dependency-free WebSocket message broker that compiles to a single binary.

It ships with its own browser client, [RosenApp](https://github.com/shivanshkc/rosenapp), a lightweight chat application powered by Rosenbridge.

## Quickstart

1. Build the project.

```bash
go build -o bin/rosenbridge cmd/rosenbridge/main.go
```

2. Create a config file. Save it anywhere you like. The example below suffices for running locally.

```json
{
  "httpServer": {
    "addr": "localhost:8080",
    "allowedOrigins": ["*"],
    "corsMaxAgeSec": 86400
  },
  "logger": {
    "level": "debug",
    "pretty": true
  },
  "database": {
    "usersFilePath": "./secrets/users.json"
  },
  "frontend": {
    "backendAddr": "http://localhost:8080",
    "path": "./client/web"
  }
}
```

3. Run Rosenbridge.

```bash
bin/rosenbridge -config <path to your config file>
```

4. Visit [http://localhost:8080](http://localhost:8080) to use RosenApp.

## API Docs

All routes are prefixed with `/api` and use Basic Auth where noted.

| Method | Path           | Auth  | Description                         |
|--------|----------------|-------|-------------------------------------|
| `GET`  | `/api`         | —     | Health check                        |
| `POST` | `/api/user`    | —     | Create a new user                   |
| `POST` | `/api/message` | Basic | Send a message to one or more users |
| `GET`  | `/api/connect` | Basic | Upgrade to WebSocket                |

**WebSocket** - Connect to `ws://<host>/api/connect` with Basic Auth. The server pushes `MessageReceived` events to the client when messages are sent to the connected user.

See [docs/API Docs.md](docs/API%20Docs.md) for full details, request/response schemas, and the middleware stack.

## Design Choices

**Why a file for a database?**  
Rosenbridge is meant to be dependency-free. It will never rely on a separate database server. In this release, a single JSON file proved sufficient. Future releases may move to an embeddable database.

**Why only Basic Auth?**  
This is the MVP release, so the feature set was kept minimal. JWT auth is planned for a future release.

## Roadmap

### Rosenbridge
1. Distributed processing to support large number of clients
2. JWT auth
3. Message delivery response
4. TCP connectivity alongside WebSocket

### RosenApp
1. Message persistence
2. Single sign-on (Google)

## License

MIT - see [LICENSE](LICENSE).
