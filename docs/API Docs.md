# API Documentation

All API routes are prefixed with `/api`. Requests and responses use JSON. Authenticated endpoints use
[HTTP Basic Auth](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#basic_authentication_scheme).

## Error Response

Every error follows this shape:

```json
{
  "status": "Bad Request",
  "reason": "human-readable explanation"
}
```

A `X-Correlation-ID` header is included in every response for tracing.

---

## `GET /api` — Health Check

Returns a simple status to confirm the server is running.

**Auth:** None

**Response — `200 OK`**

```json
{ "code": "OK" }
```

---

## `POST /api/user` — Create User

Registers a new user account.

**Auth:** None

**Request Body**

| Field      | Type   | Rules                                    |
|------------|--------|------------------------------------------|
| `username` | string | 3–100 chars, alphanumeric / `_` / `-`    |
| `password` | string | 3–100 chars                              |

**Response — `201 Created`**

```json
{ "username": "alice" }
```

**Errors**

| Status | When |
|--------|------|
| `400`  | Invalid or missing fields |
| `409`  | Username already taken |

---

## `POST /api/message` — Send Message

Sends a message to one or more connected users.

**Auth:** Basic Auth (required)

**Request Body**

| Field       | Type     | Rules                                             |
|-------------|----------|---------------------------------------------------|
| `message`   | string   | Non-empty, max 4 096 UTF-8 runes                  |
| `receivers` | string[] | 1–100 usernames, each following username rules     |

**Response — `202 Accepted`**

```json
{}
```

**Errors**

| Status | When |
|--------|------|
| `400`  | Invalid body, empty message, or bad receiver list |
| `401`  | Missing or invalid credentials |

---

## `GET /api/connect` — WebSocket Upgrade

Upgrades the HTTP connection to a WebSocket. See the [WebSocket](#websocket) section below.

**Auth:** Basic Auth — either via the `Authorization` header or query parameters `?username=<u>&password=<p>` (useful for browser clients that cannot set headers on the upgrade request).

**Response — `101 Switching Protocols`** on success.

**Errors**

| Status | When |
|--------|------|
| `401`  | Missing or invalid credentials |

---

## `GET /*` — SPA / Static Files

Serves the bundled front-end (RosenApp) from the configured `frontend.path`. Unknown paths fall back to `index.html` for client-side routing.

---

## WebSocket

### Connecting

```
ws://localhost:8080/api/connect
```

Or with TLS:

```
wss://host/api/connect
```

Authenticate with Basic Auth (header or query params as described above). Once the upgrade succeeds, the server registers the connection under the authenticated username.

### Connection Lifecycle

1. Client opens a WebSocket to `/api/connect` with credentials.
2. Server upgrades the connection and stores it by username.
3. Server runs a read loop to detect disconnects (the client does **not** send application-level events over the socket).
4. When another user calls `POST /api/message` targeting this username, the server writes a `MessageReceived` event to the socket.
5. The connection is cleaned up when the read loop exits (close frame or error).

### Server → Client Events

The server sends JSON text frames with this envelope:

```json
{
  "event_type": "<EventName>",
  "event_body": { }
}
```

#### `MessageReceived`

Delivered when a message is sent to the connected user.

```json
{
  "event_type": "MessageReceived",
  "event_body": {
    "message": "Hey!",
    "sender": "bob"
  }
}
```

| Field     | Type   | Description              |
|-----------|--------|--------------------------|
| `message` | string | The message text         |
| `sender`  | string | Username of the sender   |

### Client → Server Events

None. The client does not send application-level messages over the WebSocket. Messages are sent via the `POST /api/message` REST endpoint instead.

---

## Middleware Stack

Middleware is applied in order on every request:

| # | Middleware | Purpose |
|---|-----------|---------|
| 1 | Recovery | Catches panics; returns `500` |
| 2 | Access Logger | Logs method, URL, latency, status; generates `X-Correlation-ID` |
| 3 | CORS | Validates origins against `allowedOrigins`; handles preflight |
| 4 | Body Size Limit | Rejects request bodies larger than 16 KB |

**CORS Details**
- Allowed methods: `GET, POST, PUT, PATCH, DELETE, OPTIONS`
- Allowed headers: `Accept, Authorization, Content-Type, X-Correlation-ID`
- Exposed headers: `X-Correlation-ID`
