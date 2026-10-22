# Simple Echo Server

A minimal HTTP echo server written in Go. It echoes back details of incoming
requests as JSON and can simulate latency and custom response status codes.

## Running

```sh
just run
```

Or directly:

```sh
go run main.go
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-port` | `8000` | Port to listen on |
| `-name` | `gocho` | Server name |
| `-debug` | `false` | Enable debug logging |

## Configuration

Response behavior is controlled via `X-Echo-*` request headers:

- `X-Echo-Delay: 2` — delay the response by the given number of seconds (or a Go
  duration such as `500ms` or `2s`).
- `X-Echo-Status: 404` — respond with the given HTTP status code. Defaults to
  `200` if the value is missing or invalid.

Status values outside the valid range (`200`–`599`) or unparseable values fall
back to a `200` response.

## Examples

Echo a request (default `200`):

```sh
curl http://localhost:8000/
```

Delay the response by 2 seconds:

```sh
curl -H "X-Echo-Delay: 2" http://localhost:8000/
```

Respond with a custom status code:

```sh
curl -H "X-Echo-Status: 503" http://localhost:8000/
```

## Notes

- Query-parameter-based settings (`?delay=`, `?status=`) have been replaced by
  header-based settings (`X-Echo-*`).
- Header values are reflected back in the echoed JSON response under the
  `Headers` map.
- Responses are always JSON-encoded.
