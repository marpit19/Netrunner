# Netrunner: Building HTTP from the Ground Up
## Part 3: HTTP Responses and Status Codes

In this third part of our series on building HTTP from the ground up, we'll focus on implementing HTTP responses and status codes. We'll also create a simple server to tie everything together.

### Implementing Status Codes

Let's start by creating a new package for our HTTP status codes. First, create a new directory:

```bash
mkdir -p netrunner/pkg/http/status
```

Now, create a new file `netrunner/pkg/http/status/status.go`:

```go
package status

const (
	OK                  = 200
	Created             = 201
	Accepted            = 202
	NoContent           = 204
	MovedPermanently    = 301
	Found               = 302
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	MethodNotAllowed    = 405
	IamATeaPot          = 418
	InternalServerError = 500
	NotImplemented      = 501
	BadGateway          = 502
	ServiceUnavailable  = 503
)

var statusText = map[int]string{
	OK:                  "OK",
	Created:             "Created",
	Accepted:            "Accepted",
	NoContent:           "No Content",
	MovedPermanently:    "Moved Permanently",
	Found:               "Found",
	BadRequest:          "Bad Request",
	Unauthorized:        "Unauthorized",
	Forbidden:           "Forbidden",
	NotFound:            "Not Found",
	MethodNotAllowed:    "Method Not Allowed",
	IamATeaPot:          "I'm a teapot",
	InternalServerError: "Internal Server Error",
	NotImplemented:      "Not Implemented",
	BadGateway:          "Bad Gateway",
	ServiceUnavailable:  "Service Unavailable",
}

func Text(code int) string {
	return statusText[code]
}
```

### Updating the Response Structure

Next, let's update our `netrunner/pkg/http/response.go` file to use this new package:

```go
package http

import (
	"fmt"
	"strings"

	"github.com/yourusername/netrunner/pkg/http/status"
)

type Response struct {
	Version    string
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func NewResponse() *Response {
	return &Response{
		Version: "HTTP/1.1",
		Headers: make(map[string]string),
	}
}

func FormatResponse(r *Response) []byte {
	var builder strings.Builder

	// Status line
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Version, r.StatusCode, r.StatusText))

	// Headers
	for key, value := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Empty line
	builder.WriteString("\r\n")

	// Body
	return append([]byte(builder.String()), r.Body...)
}

func StatusText(code int) string {
	return status.Text(code)
}
```

### Implementing the Main Server

Now, let's create our `main.go` file in the `cmd/server/` directory to bring everything together:

```go
package main

import (
	"fmt"
	"io"
	"net"

	"github.com/yourusername/netrunner/pkg/http"
	"github.com/yourusername/netrunner/pkg/http/status"
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

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	request, err := http.ParseRequest(buffer[:n])
	if err != nil {
		fmt.Printf("Error parsing request: %v\n", err)
		sendErrorResponse(conn, status.BadRequest)
		return
	}

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
```

### Running and Testing the Server

To run the server:

1. Navigate to the `cmd/server/` directory in your terminal.
2. Run the following command:

```bash
go run main.go
```

You should see the message "Server listening on :8080" printed to the console.

You can test the server using a web browser or a tool like curl. For example:

```bash
curl http://localhost:8080/test
```

You should receive a response that echoes back the details of your request.

### Conclusion

In this part, we've implemented HTTP responses and status codes, and created a simple server that can handle incoming requests and send responses. This lays the foundation for building more complex HTTP functionality in future parts of this series.

In the next part, we'll expand on this implementation by adding support for different HTTP methods (GET, POST, etc.) and creating a simple routing system to handle different paths.

Stay tuned for Part 4, where we'll dive into implementing GET and POST methods and create a basic routing system!