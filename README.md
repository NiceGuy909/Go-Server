# Simple HTTP server with TCP sockets

This small example uses Go's `syscall` package instead of `net/http`.

Run it with:

```sh
go run .
```

Then open <http://localhost:8080/>. The server returns a page for `/`, a 404
for other paths, and a 405 if the method is not GET.

## How to explain it

1. `Socket`, `Bind`, and `Listen` create a TCP server on port 8080.
2. The loop calls `Accept` whenever a client connects.
3. The handler reads bytes until the blank line ending the HTTP headers.
4. It reads the method and path from the request line and chooses a response.
5. It writes an HTTP status line, headers, and body, then closes the connection.

Try a raw request with:

```sh
printf 'GET / HTTP/1.1\r\nHost: localhost\r\n\r\n' | nc localhost 8080
```

Run the tests with `go test ./...`.

This is a learning example, not a complete HTTP implementation.
