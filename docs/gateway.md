# Gateway

The **gateway** is a local WebSocket server that exposes the lingti-bot AI agent to any custom client — web UIs, mobile apps, scripts, or other tools — over a simple JSON protocol.

Instead of a specific messaging platform (Telegram, WeCom), the gateway gives you a raw API: connect a WebSocket, send chat messages, and stream back AI responses. You keep full control of the UI.

## Why use the gateway?

| Use case | Example |
|----------|---------|
| Custom web UI | React chat interface connected to your local AI |
| Mobile companion app | iOS/Android app that talks to lingti-bot |
| CLI tooling | Shell scripts that query the AI agent |
| Relay bridge | Other bots that forward messages to lingti-bot |
| Automation | Headless scripts that drive the agent programmatically |

## Starting the gateway

```bash
lingti-bot gateway --api-key sk-ant-xxx
```

The gateway listens on `:18789` by default. Change it with `--addr`:

```bash
lingti-bot gateway --addr :9000 --api-key sk-ant-xxx
```

All flags also accept environment variables — see the [CLI reference](../docs/cli-reference.md) for the full table.

## Authentication

By default the gateway accepts all connections with no auth check.

To require clients to authenticate, set one or more **auth tokens**:

```bash
# Single token
lingti-bot gateway --auth-token my-secret --api-key sk-ant-xxx

# Multiple tokens (each person gets their own token)
lingti-bot gateway --auth-tokens "alice-token,bob-token" --api-key sk-ant-xxx

# Via environment variable
GATEWAY_AUTH_TOKENS="alice-token,bob-token" lingti-bot gateway --api-key sk-ant-xxx
```

When auth is enabled, a client **must** send an `auth` message before it can chat. Any token in the list grants full access — there is no per-token permission hierarchy.

## HTTP endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Returns `{"status":"ok"}` — useful for health checks |
| `GET` | `/status` | Returns running status, client count, address, and whether auth is enabled |
| `GET` | `/ws` | WebSocket upgrade endpoint |

```bash
# Health check
curl http://localhost:18789/health

# Status
curl http://localhost:18789/status
# {"addr":":18789","auth_enabled":true,"clients":2,"status":"running"}
```

## WebSocket protocol

All messages are JSON objects with this envelope:

```json
{
  "id":        "unique-message-id",
  "type":      "message-type",
  "payload":   { ... },
  "timestamp": 1700000000000
}
```

- `id` — arbitrary string; the server echoes it back on response messages so you can match requests to responses
- `type` — one of the types listed below
- `payload` — type-specific object (omitted for ping/pong)
- `timestamp` — milliseconds since epoch (set by sender; server always sets it on outbound messages)

### Client → Server message types

#### `ping`

Keep-alive from the client. Server replies with `pong`.

```json
{"type": "ping"}
```

The server also sends WebSocket-level `Ping` frames every 30 seconds; respond with `Pong` to keep the connection alive.

#### `auth`

Authenticate with a token. Required before `chat`/`command` when auth is enabled.

```json
{
  "type": "auth",
  "payload": {
    "token": "my-secret-token"
  }
}
```

Server replies with `auth_result`.

#### `chat`

Send a message to the AI agent.

```json
{
  "id": "req-1",
  "type": "chat",
  "payload": {
    "text": "What is the capital of France?",
    "session_id": "optional-session-id"
  }
}
```

- `text` — the message to send (required)
- `session_id` — reuse an existing conversation session; omit to continue the current session for this connection, or to start fresh

Server replies with one or more `response` messages.

#### `command`

Send a built-in gateway command.

```json
{
  "type": "command",
  "payload": {
    "command": "status"
  }
}
```

**Available commands:**

| Command | Description |
|---------|-------------|
| `status` | Returns this connection's `client_id`, `session_id`, and `authorized` state |
| `clear` | Clears the session ID so the next message starts a new conversation |

### Server → Client message types

#### `pong`

Reply to a client `ping`.

```json
{"type": "pong", "timestamp": 1700000000000}
```

#### `auth_result`

Result of an `auth` attempt.

```json
{
  "type": "auth_result",
  "payload": {
    "success": true,
    "message": ""
  }
}
```

On failure, `success` is `false` and `message` explains why (e.g. `"Invalid token"`).

#### `response`

An AI response chunk. The server sends one response per `chat` message (currently non-streaming — the full text arrives at once with `done: true`).

```json
{
  "id": "req-1",
  "type": "response",
  "payload": {
    "text": "The capital of France is Paris.",
    "session_id": "abc123",
    "done": true
  }
}
```

- `id` echoes the request `id` from the `chat` message
- `done: true` means this is the final response for that request

#### `event`

A server-pushed event, in response to a `command`.

```json
{
  "type": "event",
  "payload": {
    "event": "status",
    "data": {
      "client_id": "1700000000000",
      "session_id": "abc123",
      "authorized": true
    }
  }
}
```

#### `error`

An error occurred processing the last message.

```json
{
  "type": "error",
  "payload": {
    "code": "unauthorized",
    "message": "Authentication required"
  }
}
```

**Error codes:**

| Code | Meaning |
|------|---------|
| `unauthorized` | Auth is enabled and the client has not authenticated |
| `invalid_message` | Message could not be parsed as JSON |
| `invalid_payload` | Payload fields are missing or malformed |
| `handler_error` | The AI agent returned an error |
| `no_handler` | Gateway has no message handler (internal error) |
| `unknown_type` | Unknown message `type` field |
| `unknown_command` | Unknown `command` value in a `command` message |

## Sessions

Each connection starts with no session. The first `chat` message creates a session:

1. If `session_id` is provided in the payload, that ID is used.
2. Otherwise the connection's own client ID becomes the session ID.

Subsequent `chat` messages on the same connection reuse that session, so the AI remembers earlier context. To start fresh, send a `clear` command or provide a new `session_id`.

## Example: JavaScript client

```js
const ws = new WebSocket("ws://localhost:18789/ws");

ws.onopen = () => {
  // Authenticate (skip this step if no auth tokens are configured)
  ws.send(JSON.stringify({
    type: "auth",
    payload: { token: "my-secret-token" }
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  if (msg.type === "auth_result") {
    if (!msg.payload.success) {
      console.error("Auth failed:", msg.payload.message);
      return;
    }
    // Authenticated — send a chat message
    ws.send(JSON.stringify({
      id: "req-1",
      type: "chat",
      payload: { text: "Hello, what can you do?" }
    }));
  }

  if (msg.type === "response" && msg.payload.done) {
    console.log("AI:", msg.payload.text);
  }

  if (msg.type === "error") {
    console.error(`Error [${msg.payload.code}]:`, msg.payload.message);
  }
};
```

## Example: Python client

```python
import json
import websocket  # pip install websocket-client

ws = websocket.create_connection("ws://localhost:18789/ws")

# Authenticate
ws.send(json.dumps({"type": "auth", "payload": {"token": "my-secret-token"}}))
result = json.loads(ws.recv())
assert result["payload"]["success"], f"Auth failed: {result['payload']['message']}"

# Send a message
ws.send(json.dumps({
    "id": "req-1",
    "type": "chat",
    "payload": {"text": "Summarize today's news"}
}))

# Read response
response = json.loads(ws.recv())
print("AI:", response["payload"]["text"])

ws.close()
```

## Connection lifecycle

```
Client                          Gateway
  |                                |
  |--- WebSocket upgrade --------->|
  |<-- connection accepted --------|
  |                                |
  |--- auth (if required) -------->|
  |<-- auth_result ----------------|
  |                                |
  |--- chat {"text": "Hi"} ------->|
  |<-- response {"done": true} ----|
  |                                |
  |--- command {"command":"clear"} |
  |<-- event {"event":"cleared"} --|
  |                                |
  |--- [disconnect] -------------->|
```

The server sends WebSocket-level Ping frames every 30 seconds. If your client does not respond with Pong within 60 seconds, the connection is closed.
