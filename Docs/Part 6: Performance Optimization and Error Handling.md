# Netrunner: Building HTTP from the Ground Up
## Part 6: Performance Optimization and Error Handling

Welcome to Part 6 of our series on building HTTP from the ground up! In this installment, we'll focus on improving the performance of our server and implementing robust error handling.

### 1. Connection Pooling

One way to improve performance is by implementing connection pooling. This allows us to reuse connections instead of creating a new one for each request.

Let's create a new file `pkg/http/connpool.go`:

```go
package http

import (
    "net"
    "sync"
    "time"
)

type ConnPool struct {
    mu       sync.Mutex
    conns    chan net.Conn
    maxConns int
}

func NewConnPool(maxConns int) *ConnPool {
    return &ConnPool{
        conns:    make(chan net.Conn, maxConns),
        maxConns: maxConns,
    }
}

func (p *ConnPool) Get(network, address string) (net.Conn, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    select {
    case conn := <-p.conns:
        return conn, nil
    default:
        return net.Dial(network, address)
    }
}

func (p *ConnPool) Put(conn net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()

    select {
    case p.conns <- conn:
    default:
        conn.Close()
    }
}

func (p *ConnPool) CloseIdleConnections() {
    p.mu.Lock()
    defer p.mu.Unlock()

    close(p.conns)
    for conn := range p.conns {
        conn.Close()
    }
    p.conns = make(chan net.Conn, p.maxConns)
}
```

Now, let's update our `main.go` to use this connection pool:

```go
// In main.go

var connPool *http.ConnPool

func main() {
    // ... (previous code)

    connPool = http.NewConnPool(100)
    
    // ... (rest of the code)
}

func handleConnection(conn net.Conn, router *http.Router) {
    defer connPool.Put(conn)

    // ... (rest of the function remains the same)
}
```

### 2. Timeout Handling

Let's implement timeouts to prevent long-running requests from blocking our server. Update the `handleConnection` function in `main.go`:

```go
func handleConnection(conn net.Conn, router *http.Router) {
    defer connPool.Put(conn)

    conn.SetDeadline(time.Now().Add(30 * time.Second))

    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            sendErrorResponse(conn, http.StatusRequestTimeout)
        } else {
            fmt.Printf("Error reading from connection: %v\n", err)
        }
        return
    }

    // ... (rest of the function remains the same)
}
```

### 3. Graceful Shutdown

Implement a graceful shutdown to ensure all ongoing requests are completed before the server stops. Update `main.go`:

```go
import (
    // ... (other imports)
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // ... (previous code)

    go startServer(":8080", router)

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    fmt.Println("Server is shutting down...")
    connPool.CloseIdleConnections()
    fmt.Println("Server stopped")
}
```

### 4. Error Handling

Let's improve our error handling by creating custom error types. Create a new file `pkg/http/errors.go`:

```go
package http

import "fmt"

type HTTPError struct {
    Code    int
    Message string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP error %d: %s", e.Code, e.Message)
}

func NewHTTPError(code int, message string) *HTTPError {
    return &HTTPError{Code: code, Message: message}
}
```

Now, let's update our `handleConnection` function to use these custom errors:

```go
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
        handleHTTPError(conn, http.NewHTTPError(http.StatusRequestTimeout, "Request timeout"))
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
```

### Conclusion

In this part, we've significantly improved our HTTP server's performance and error handling:

1. We implemented a connection pool to reuse connections and reduce the overhead of creating new ones.
2. We added timeout handling to prevent long-running requests from blocking the server.
3. We implemented a graceful shutdown mechanism to ensure clean server stops.
4. We improved our error handling with custom error types and more robust error responses.

These improvements make our server more resilient and better equipped to handle real-world scenarios. In the next part, we'll focus on implementing HTTPS and TLS support to make our server secure.

Stay tuned for Part 7, where we'll dive into the world of secure communications!