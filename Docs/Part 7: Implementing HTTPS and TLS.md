# Netrunner: Building HTTP from the Ground Up
## Part 7: Implementing HTTPS

Welcome to Part 7 of our series on building HTTP from the ground up! In this installment, we'll focus on implementing HTTPS (HTTP Secure) to encrypt communications between clients and our server.

### Understanding HTTPS

HTTPS is an extension of HTTP that uses TLS (Transport Layer Security) for secure communication over a computer network. It provides:

1. Encryption: Protecting the exchanged data from eavesdropping and tampering
2. Authentication: Ensuring that the server is who it claims to be
3. Integrity: Verifying that the data hasn't been forged or tampered with

### Generating a Self-Signed Certificate

For development purposes, we'll create a self-signed certificate. In a production environment, you'd use a certificate signed by a trusted Certificate Authority.

Run the following command to generate a self-signed certificate and private key:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

This will create two files: `key.pem` (private key) and `cert.pem` (certificate).

### Implementing HTTPS Server

Let's update our `main.go` to support both HTTP and HTTPS:

```go
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/netrunner/pkg/http"
)

var connPool *http.ConnPool

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

	connPool = http.NewConnPool(100)

	// Start HTTP server
	go startServer("http", ":8080", router)

	// Start HTTPS server
	go startServer("https", ":8443", router)

	// Wait for interrupt signal to gracefully shut down the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Server is shutting down...")
	connPool.CloseIdleConnections()
	fmt.Println("Server stopped")
}

func startServer(protocol string, address string, router *http.Router) {
	var listener net.Listener
	var err error

	if protocol == "https" {
		// Load TLS certificate and key
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			fmt.Printf("Failed to load TLS certificate: %v\n", err)
			return
		}

		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		listener, err = tls.Listen("tcp", address, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", address)
	}

	if err != nil {
		fmt.Printf("Failed to start %s server: %v\n", protocol, err)
		return
	}
	defer listener.Close()

	fmt.Printf("%s server listening on %s\n", protocol, address)

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
	defer connPool.Put(conn)

	var tlsConn *tls.ConnectionState
	if tlsConnection, ok := conn.(*tls.Conn); ok {
		state := tlsConnection.ConnectionState()
		tlsConn = &state
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		handleConnectionError(conn, err)
		return
	}

	request, err := http.ParseRequest(buffer[:n], tlsConn)
	if err != nil {
		handleHTTPError(conn, http.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err)))
		return
	}

	response := router.HandleRequest(request)
	_, err = conn.Write(http.FormatResponse(response))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

// ... (other handler functions remain the same)
```

### Updating the Request Structure

We need to update our `Request` structure to include TLS information. Update `pkg/http/request.go`:

```go
package http

import "crypto/tls"

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
	TLS     *tls.ConnectionState
}

// Update ParseRequest function to set TLS info
func ParseRequest(data []byte, tlsConn *tls.ConnectionState) (*Request, error) {
	// ... (previous parsing logic)

	request.TLS = tlsConn

	return request, nil
}
```

### Testing HTTPS

To test HTTPS functionality:

1. Run your server:
   ```bash
   go run cmd/server/main.go
   ```

2. Open a web browser and navigate to `https://localhost:8443`
   (You'll see a security warning because we're using a self-signed certificate. In a real-world scenario, you'd use a certificate from a trusted CA.)

3. You can also use curl to test:
   ```bash
   curl -k https://localhost:8443
   ```
   The `-k` option tells curl to accept self-signed certificates.

### Conclusion

In this part, we've successfully implemented HTTPS support for our Netrunner server:

1. We generated a self-signed certificate for testing purposes.
2. We updated our server to support both HTTP and HTTPS connections.
3. We modified our `Request` structure to include TLS information.

These changes significantly improve the security of our server by encrypting all communications over HTTPS.

If you want to get new part pls follow me on twitter: [my twitter](https://x.com/minamisatokun) to stay tuned