# Netrunner: Building HTTP from the Ground Up
## Part 1: Building a Basic TCP Server

Welcome to Part 1 of our Netrunner series! In this part, we'll implement a basic TCP (Transmission Control Protocol) server using Go. This will serve as the foundation for our HTTP server in later parts.

### Understanding TCP

Before we dive into coding, let's briefly review what TCP is:

- TCP is a connection-oriented protocol that provides reliable, ordered, and error-checked delivery of data between applications running on hosts communicating over an IP network.
- It's the underlying protocol for many application-layer protocols, including HTTP.
- TCP uses a three-way handshake to establish a connection and ensures data integrity through sequence numbers and acknowledgments.

### Setting Up Our Project

First, let's set up our project structure:

```bash
mkdir -p netrunner/pkg/tcp
touch netrunner/pkg/tcp/server.go
mkdir -p netrunner/test/tcp
touch netrunner/test/tcp/server_test.go
```

### Writing Our First Test

Following Test-Driven Development (TDD) principles, we'll start by writing a test. Open `netrunner/test/tcp/server_test.go` and add the following:

```go
package tcp_test

import (
	"net"
	"testing"
	"time"

	"github.com/yourusername/netrunner/pkg/tcp"
)

func TestTCPServer(t *testing.T) {
	// Start the server
	go tcp.StartServer("localhost:8080")

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	t.Log("Successfully connected to the server")
}
```

This test attempts to start our TCP server and then connects to it. If the connection is successful, we know our server is working.

### Implementing the TCP Server

Now, let's implement the `StartServer` function in `netrunner/pkg/tcp/server.go`:

```go
package tcp

import (
	"fmt"
	"net"
)

// StartServer initializes and starts a TCP server on the given address
func StartServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", address)

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
	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
	// We'll implement actual handling in the next part
}
```

Let's break down what this code does:

1. `net.Listen("tcp", address)` creates a TCP listener on the specified address.
2. We enter an infinite loop to continuously accept new connections.
3. For each new connection, we spawn a goroutine to handle it, allowing our server to handle multiple connections concurrently.

### Running the Test

To run the test, navigate to your project root and run:

```bash
go test ./test/tcp
```

If everything is set up correctly, the test should pass.

### Adding Echo Functionality

Let's enhance our server by adding simple echo functionality. Update the `handleConnection` function in `server.go`:

```go
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			return
		}

		fmt.Printf("Received: %s", string(buffer[:n]))
		
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}
```

This implementation reads data from the connection and writes it back, creating an echo server.

### Testing Echo Functionality

Add a new test to `server_test.go`:

```go
func TestEchoFunctionality(t *testing.T) {
	go tcp.StartServer("localhost:8081")
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	testMessage := "Hello, Netrunner!"
	_, err = conn.Write([]byte(testMessage))
	if err != nil {
		t.Fatalf("Could not send message to server: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Could not read from server: %v", err)
	}

	response := string(buffer[:n])
	if response != testMessage {
		t.Fatalf("Expected response %q, but got %q", testMessage, response)
	}

	t.Log("Echo functionality working correctly")
}
```

Run the tests again to ensure everything is working as expected.

### Conclusion

We've now implemented a basic TCP server with echo functionality, following TDD principles. This server forms the foundation for our HTTP server. Key takeaways from this part:

1. We used Go's `net` package to create a TCP listener and handle connections.
2. We implemented concurrent handling of multiple connections using goroutines.
3. We added basic echo functionality to demonstrate reading from and writing to connections.

In the next part, we'll start implementing HTTP-specific functionality on top of this TCP server. We'll learn about HTTP request structure and begin parsing incoming HTTP requests.

Remember, this is a simplified implementation for educational purposes. Production servers need to handle various edge cases, implement proper error handling, and often include features like connection pooling and timeouts.

Stay tuned for Part 2, where we'll dive into HTTP basics and request parsing!