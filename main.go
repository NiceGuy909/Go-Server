package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"syscall"
)

func main() {
	server, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(server)

	address := &syscall.SockaddrInet4{Port: 8080}
	if err := syscall.Bind(server, address); err != nil {
		log.Fatal(err)
	}
	if err := syscall.Listen(server, 10); err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on http://localhost:8080")
	for {
		client, _, err := syscall.Accept(server)
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleConnection(client)
	}
}

func handleConnection(client int) {
	defer syscall.Close(client)

	request := make([]byte, 0, 1024)
	buffer := make([]byte, 1024)
	for !bytes.Contains(request, []byte("\r\n\r\n")) {
		n, err := syscall.Read(client, buffer)
		if err != nil || n == 0 {
			writeResponse(client, "400 Bad Request", "bad request\n")
			return
		}
		request = append(request, buffer[:n]...)
		if len(request) > 8192 {
			writeResponse(client, "400 Bad Request", "request too large\n")
			return
		}
	}

	lines := strings.SplitN(string(request), "\r\n", 2)
	fields := strings.Fields(lines[0])
	if len(fields) != 3 {
		writeResponse(client, "400 Bad Request", "bad request\n")
		return
	}

	method, path := fields[0], fields[1]
	switch {
	case method != "GET":
		writeResponse(client, "405 Method Not Allowed", "method not allowed\n")
	case path == "/":
		writeResponse(client, "200 OK", "<h1>Hello from Go!</h1>\n")
	default:
		writeResponse(client, "404 Not Found", "page not found\n")
	}
}

func writeResponse(client int, status, body string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Type: text/html; charset=utf-8\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
		status, len(body), body,
	)
	if _, err := syscall.Write(client, []byte(response)); err != nil {
		log.Println("write:", err)
	}
}
