# Netrunner: Building HTTP from the Ground Up
## Part 4: HTTP Methods and Routing

Welcome to Part 4 of our series on building HTTP from the ground up! In this installment, we'll focus on implementing GET and POST methods and creating a basic routing system for our server.

### Implementing HTTP Methods

First, let's modify our `ParseRequest` function to handle the request body:

```go
func ParseRequest(data []byte) (*Request, error) {
    // ... (previous code for parsing the request line and headers)

    // Parse the body
    bodyStart := bytes.Index(data, []byte("\r\n\r\n")) + 4
    if bodyStart < len(data) {
        request.Body = data[bodyStart:]
    }

    return request, nil
}
```

### Creating a Basic Router

Let's create a new file `pkg/http/router.go`:

```go
package http

import (
    "fmt"
    "strings"
)

type HandlerFunc func(*Request) *Response

type Router struct {
    routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
    return &Router{
        routes: make(map[string]map[string]HandlerFunc),
    }
}

func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
    if _, ok := r.routes[method]; !ok {
        r.routes[method] = make(map[string]HandlerFunc)
    }
    r.routes[method][path] = handler
}

func (r *Router) HandleRequest(req *Request) *Response {
    if handlers, ok := r.routes[req.Method]; ok {
        if handler, ok := handlers[req.Path]; ok {
            return handler(req)
        }
    }
    return NotFoundResponse()
}

func NotFoundResponse() *Response {
    resp := NewResponse()
    resp.StatusCode = 404
    resp.StatusText = "Not Found"
    resp.SetBody([]byte("404 - Not Found"))
    return resp
}
```

### Updating the Main Server

Now, let's update our `main.go` to use the new router:

```go
package main

import (
	"fmt"
	"io"
	"net"

	"github.com/appyzdl/Netrunner/pkg/http"
	"github.com/appyzdl/Netrunner/pkg/http/status"
)

func main() {
	router := http.NewRouter()

	// Add routes
	router.AddRoute("GET", "/", handleRoot)
	router.AddRoute("GET", "/hello", handleHello)
	router.AddRoute("POST", "/echo", handleEcho)

	startServer(":8080", router)
}

func startServer(address string, router *http.Router) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConnection(conn, router)
	}
}

func handleConnection(conn net.Conn, router *http.Router) {
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

	response := router.HandleRequest(request)
	_, err = conn.Write(http.FormatResponse(response))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

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

```

### Testing the Server

Now you can run the server and test it with curl:

1. Start the server:
   ```
   go run main.go
   ```

2. Test GET requests:
   ```
   curl http://localhost:8080/
   curl http://localhost:8080/hello
   ```

3. Test POST request:
   ```
   curl -X POST -d "Hello, Netrunner!" http://localhost:8080/echo
   ```

4. Test non-existent route:
   ```
   curl http://localhost:8080/notfound
   ```

### Conclusion

Please refer to code there might be some extra changes, thank you!!

In this part, we've implemented a basic routing system and added support for GET and POST methods. Our Netrunner server can now handle different routes and HTTP methods, laying the groundwork for more complex web applications.

In the next part, we'll focus on improving our server by adding support for static file serving and implementing basic middleware functionality. Stay tuned for Part 5, where we'll continue to enhance our HTTP server!