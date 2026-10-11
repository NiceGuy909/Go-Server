package main

import (
	"strings"
	"syscall"
	"testing"
)

func TestResponses(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		status string
		body   string
	}{
		{"home page", "GET", "/", "200 OK", "<h1>Hello from Go!</h1>\n"},
		{"unknown page", "GET", "/missing", "404 Not Found", "page not found\n"},
		{"unsupported method", "POST", "/", "405 Method Not Allowed", "method not allowed\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := test.method + " " + test.path + " HTTP/1.1\r\nHost: localhost\r\n\r\n"
			response := sendRequest(t, request)
			header, body, found := strings.Cut(response, "\r\n\r\n")
			if !found {
				t.Fatalf("invalid response: %q", response)
			}
			if !strings.Contains(header, "HTTP/1.1 "+test.status) {
				t.Errorf("unexpected status line: %q", header)
			}
			if body != test.body {
				t.Errorf("body = %q, want %q", body, test.body)
			}
		})
	}
}

func sendRequest(t *testing.T, request string) string {
	t.Helper()

	sockets, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(sockets[1])

	go handleConnection(sockets[0])
	if _, err := syscall.Write(sockets[1], []byte(request)); err != nil {
		t.Fatal(err)
	}

	var response []byte
	buffer := make([]byte, 1024)
	for {
		n, err := syscall.Read(sockets[1], buffer)
		if err != nil {
			t.Fatal(err)
		}
		if n == 0 {
			return string(response)
		}
		response = append(response, buffer[:n]...)
	}
}
