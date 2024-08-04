# Netrunner: Building HTTP from the Ground Up
## Part 5: Middleware and Static File Serving

Welcome to Part 5 of our series on building HTTP from the ground up! In this installment, we'll enhance our server by implementing middleware functionality and adding support for serving static files.

### Implementing Middleware

Middleware allows us to add reusable components to our request handling pipeline. Let's start by updating our router to support middleware.

First, let's update our `pkg/http/router.go` file:

```go
package http

type HandlerFunc func(*Request) *Response
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Router struct {
    routes     map[string]map[string]HandlerFunc
    middleware []MiddlewareFunc
}

func NewRouter() *Router {
    return &Router{
        routes:     make(map[string]map[string]HandlerFunc),
        middleware: []MiddlewareFunc{},
    }
}

func (r *Router) Use(mw MiddlewareFunc) {
    r.middleware = append(r.middleware, mw)
}

func (r *Router) HandleRequest(req *Request) *Response {
    if handlers, ok := r.routes[req.Method]; ok {
        if handler, ok := handlers[req.Path]; ok {
            // Apply middleware
            for i := len(r.middleware) - 1; i >= 0; i-- {
                handler = r.middleware[i](handler)
            }
            return handler(req)
        }
    }
    return NotFoundResponse()
}
```

Now let's implement a simple logging middleware. Create a new file `pkg/http/middleware.go`:

```go
package http

import (
    "fmt"
    "time"
)

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
    return func(req *Request) *Response {
        start := time.Now()
        resp := next(req)
        duration := time.Since(start)
        fmt.Printf("%s %s - %d (%v)\n", req.Method, req.Path, resp.StatusCode, duration)
        return resp
    }
}
```

### Serving Static Files

Next, let's add support for serving static files. We'll create a new handler function for this in `pkg/http/handlers.go`:

```go
package http

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func StaticFileHandler(basePath string) HandlerFunc {
    return func(req *Request) *Response {
        filePath := strings.TrimPrefix(req.Path, "/static")
        
        if filePath == "" || filePath == "/" {
            filePath = "/index.html"
        }
        
        if strings.Contains(filePath, "..") {
            return NotFoundResponse()
        }

        fullPath := filepath.Join(basePath, filePath)
        
        fmt.Printf("Attempting to serve file: %s\n", fullPath) // Debug log

        file, err := os.Open(fullPath)
        if err != nil {
            fmt.Printf("Error opening file: %v\n", err) // Debug log
            return NotFoundResponse()
        }
        defer file.Close()

        stat, err := file.Stat()
        if err != nil {
            fmt.Printf("Error getting file stats: %v\n", err) // Debug log
            return InternalServerErrorResponse()
        }

        if stat.IsDir() {
            fmt.Println("Requested path is a directory") // Debug log
            return NotFoundResponse()
        }

        content, err := os.ReadFile(fullPath)
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err) // Debug log
            return InternalServerErrorResponse()
        }

        resp := NewResponse()
        resp.StatusCode = StatusOK
        resp.StatusText = StatusText(StatusOK)
        resp.SetHeader("Content-Type", getContentType(fullPath))
        resp.SetHeader("Content-Length", fmt.Sprintf("%d", len(content)))
        resp.Body = content
        return resp
    }
}

func getContentType(path string) string {
    ext := filepath.Ext(path)
    
    if mimeType := mime.TypeByExtension(ext); mimeType != "" {
        return mimeType
    }
    
    switch ext {
    case ".html", ".htm":
        return "text/html"
    case ".css":
        return "text/css"
    case ".js":
        return "application/javascript"
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".gif":
        return "image/gif"
    case ".svg":
        return "image/svg+xml"
    case ".xml":
        return "application/xml"
    case ".txt":
        return "text/plain"
    case ".pdf":
        return "application/pdf"
    case ".zip":
        return "application/zip"
    case ".mp3":
        return "audio/mpeg"
    case ".mp4":
        return "video/mp4"
    default:
        return "application/octet-stream"
    }
}
```

### Updating the Main Server

Now let's update our `main.go` to use these new features:

```go
package main

import (
    "fmt"
    "net"
    "path/filepath"
    "os"

    "github.com/yourusername/netrunner/pkg/http"
)

func main() {
    router := http.NewRouter()

    // Add middleware
    router.Use(http.LoggingMiddleware)

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

    fmt.Printf("Serving static files from: %s\n", publicPath) // Debug log

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
    if err != nil {
        fmt.Printf("Error reading from connection: %v\n", err)
        return
    }

    request, err := http.ParseRequest(buffer[:n])
    if err != nil {
        fmt.Printf("Error parsing request: %v\n", err)
        sendErrorResponse(conn, http.StatusBadRequest)
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

// ... (other handler functions remain the same)
```

### Testing the New Features

To test these new features:

1. Create a `public` directory in your project root.
2. Add some static files (e.g., HTML, CSS, JS) to this directory.
3. Run your server:
   ```bash
   go build -o netrunner cmd/server/main.go
   ./netrunner
   ```
4. Try accessing static files through the `/static/` route:
   ```
   http://localhost:8080/static/index.html
   ```
5. Observe the logging output for each request in your server console.

### Conclusion

In this part, we've significantly enhanced our HTTP server by adding middleware support and static file serving capabilities. These features bring our server closer to production-ready status and provide a foundation for building more complex web applications.

Key takeaways from this part:
1. We implemented a simple middleware system that allows for request/response processing.
2. We added a logging middleware to demonstrate the middleware functionality.
3. We created a static file handler that can serve files from a specified directory.
4. We implemented content type detection for various file types.
5. We updated our main server to use these new features.

In the next part, we'll focus on improving the performance and robustness of our server, including implementing connection pooling and adding proper error handling.

Stay tuned for Part 6, where we'll dive into advanced server optimizations!