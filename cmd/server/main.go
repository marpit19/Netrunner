package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/appyzdl/Netrunner/pkg/http"
	"github.com/appyzdl/Netrunner/pkg/http/status"
)

var connPool *http.ConnPool

func main() {
	router := http.NewRouter()

	// Add middleware
	router.Use(http.LoggingMiddleware)

	connPool = http.NewConnPool(100) // Create a pool with 100 max connections

	// Add routes
	router.AddRoute("GET", "/", handleRoot)
	router.AddRoute("GET", "/hello", handleHello)
	router.AddRoute("POST", "/echo", handleEcho)

	// Add static file handler
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	publicPath := filepath.Join(execDir, "public")
	staticHandler := http.StaticFileHandler(publicPath)
	router.AddRoute("GET", "/static/", staticHandler)

	// fmt.Printf("Serving static files from: %s\n", publicPath) // Debug log

	go startServer(":8080", router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGABRT)
	<-quit

	fmt.Println("Server is shutting down...🪦")
	connPool.CloseIdleConnections()
	fmt.Println("Server stopped")
}

func startServer(address string, router *http.Router) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to start server: %v 😭\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s 🙋‍♀️\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v 😔\n", err)
			continue
		}
		go handleConnection(conn, router)
	}
}

func handleConnection(conn net.Conn, router *http.Router) {
	defer connPool.Put(conn) // Return the connection to the pool

	// Set a timeout for the entire request handling
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		handleConnectionError(conn, err)
		return
	}

	request, err := http.ParseRequest(buffer[:n])
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
		handleHTTPError(conn, http.NewHTTPError(status.BadRequest, "Invalid request"))
		return
	}

	response := router.HandleRequest(request)
	_, err = conn.Write(http.FormatResponse(response))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

func handleConnectionError(conn net.Conn, err error) {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		handleHTTPError(conn, http.NewHTTPError(status.StatusRequestTimeout, "Request timeout"))
	} else {
		fmt.Printf("Error reading from connection: %v\n", err)
	}
}

func handleHTTPError(conn net.Conn, err *http.HTTPError) {
	response := http.NewResponse()
	response.StatusCode = err.Code
	response.StatusText = http.StatusText(err.Code)
	response.SetBody([]byte(err.Message))

	_, writeErr := conn.Write(http.FormatResponse(response))
	if writeErr != nil {
		fmt.Printf("Error writing error response: %v\n", writeErr)
	}
}

/*
func sendErrorResponse(conn net.Conn, statusCode int) {
	response := http.NewResponse()
	response.StatusCode = statusCode
	response.StatusText = http.StatusText(statusCode)
	response.SetBody([]byte(http.StatusText(statusCode)))

	_, err := conn.Write(http.FormatResponse(response))
	if err != nil {
		fmt.Printf("Error writing error response: %v\n", err)
	}
}
*/

// Handler functions
func handleRoot(req *http.Request) *http.Response {
	resp := http.NewResponse()
	resp.StatusCode = 200
	resp.StatusText = "OK"
	resp.SetBody([]byte("Welcome to Netrunner!"))
	return resp
}

func handleHello(req *http.Request) *http.Response {
	resp := http.NewResponse()
	resp.StatusCode = 200
	resp.StatusText = "OK"
	resp.SetBody([]byte("Hello, Netrunner!"))
	return resp
}

func handleEcho(req *http.Request) *http.Response {
	resp := http.NewResponse()
	resp.StatusCode = 200
	resp.StatusText = "OK"
	resp.SetHeader("Content-Type", "text/plain")
	resp.SetHeader("Content-Length", fmt.Sprintf("%d", len(req.Body)))
	resp.Body = req.Body
	return resp
}
