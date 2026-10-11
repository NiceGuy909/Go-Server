# Simple HTTP server with TCP sockets

This small example uses Go's `syscall` package instead of `net/http`.

Run it with:

```sh
go run .
```

Then open <http://localhost:8080/>. The server returns a page for `/`, a 404
for other paths, and a 405 if the method is not GET.

## How it works

1. **Create the listening socket.** `Socket` asks the operating system for an
   IPv4 TCP socket. `Bind` assigns it port 8080, and `Listen` marks it ready to
   receive connections.
2. **Accept a client.** The server calls `Accept` in a loop. Each successful
   call returns a new socket for that client's connection. The server handles
   each connection in a goroutine so it can continue accepting others.
3. **Read the HTTP request.** `handleConnection` reads bytes from the client
   until it sees `\r\n\r\n`, the blank line that marks the end of the headers.
   The request's first line looks like `GET / HTTP/1.1`.
4. **Choose a response.** The server splits that first line into its method,
   path, and HTTP version. It returns `200 OK` for a GET request to `/`,
   `404 Not Found` for other paths, and `405 Method Not Allowed` for methods
   other than GET.
5. **Write and close.** `writeResponse` sends the HTTP status line, a few
   headers, a blank line, and the response body. `Content-Length` tells the
   client how many body bytes to expect. The server closes the client socket
   after sending one response.

Try a raw request with:

```sh
printf 'GET / HTTP/1.1\r\nHost: localhost\r\n\r\n' | nc localhost 8080
```

Run the tests with `go test ./...`.

This is a learning example, not a complete HTTP implementation. It handles
one request per connection and does not support request bodies or every HTTP
protocol feature.
