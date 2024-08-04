package main

import (
	"fmt"
	"io"
	"net"

	"github.com/appyzdl/Netrunner/pkg/http"
	"github.com/appyzdl/Netrunner/pkg/http/status"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read the incoming HTTP request
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	// Parse the request (we'll assume it's valid for now)
	request, err := http.ParseRequest(buffer[:n])
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
		sendErrorResponse(conn, status.BadRequest)
		return
	}

	// For now, we'll just echo back some information about the request
	response := http.NewResponse()
	response.SetStatus(status.OK)
	response.SetHeader("Content-Type", "text/plain")
	responseBody := fmt.Sprintf("Received request:\nMethod: %s\nPath: %s\nProtocol: %s\n",
		request.Method, request.Path, request.Version)
	response.SetBody([]byte(responseBody))

	_, err = conn.Write(response.Write())
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

func sendErrorResponse(conn net.Conn, statusCode int) {
	response := http.NewResponse()
	response.SetStatus(statusCode)
	response.SetHeader("Content-Type", "text/plain")
	response.SetBody([]byte(status.Text(statusCode)))

	_, err := conn.Write(response.Write())
	if err != nil {
		fmt.Printf("Error writing error response: %v\n", err)
	}
}
