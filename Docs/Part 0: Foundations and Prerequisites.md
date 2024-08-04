# Netrunner: Building HTTP from the Ground Up
## Part 0: Foundations and Prerequisites

Welcome to Netrunner, an in-depth journey into building HTTP and HTTPS from the ground up using Go. This series is designed to take you from the basics of networking to implementing a fully functional web server. Let's start by covering the prerequisites and fundamental concepts.

### Prerequisites

1. **Go Programming Language**: You should have Go (version 1.16 or later) installed on your system. If not, visit [golang.org](https://golang.org) for installation instructions.

2. **Basic Go Knowledge**: Familiarity with Go syntax, types, and core concepts like goroutines and channels will be helpful.

3. **Command Line Interface**: Basic comfort with using terminal or command prompt.

4. **Text Editor or IDE**: Any text editor will do, but an IDE with Go support (like VSCode with Go extension) can be beneficial.

### Setting Up Your Environment

1. Install Go from [golang.org](https://golang.org)
2. Verify installation by running `go version` in your terminal
3. Set up your project:

```bash
mkdir -p netrunner/{cmd,pkg,internal,test}
cd netrunner
go mod init github.com/yourusername/netrunner
```

### Fundamental Concepts

#### 1. What is HTTP?

HTTP (Hypertext Transfer Protocol) is an application-layer protocol for transmitting hypermedia documents, such as HTML. It follows a client-server model, where web browsers, for example, act as clients, and web servers host the data and respond to client requests.

Key features of HTTP:
- Stateless protocol (each request is independent)
- Allows for client-server communication
- Supports various methods (GET, POST, etc.) for different types of requests

#### 2. What is HTTPS?

HTTPS (HTTP Secure) is an extension of HTTP. It uses TLS (Transport Layer Security) or, formerly, SSL (Secure Sockets Layer) for secure communication over a computer network.

Key features of HTTPS:
- Encrypted communication
- Data integrity
- Authentication

#### 3. What are Sockets?

A socket is one endpoint of a two-way communication link between two programs running on the network. It's bound to a port number so that the TCP layer can identify the application that data is destined to be sent to.

In Go, the `net` package provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets.

#### 4. How Web Communication Works

1. **DNS Resolution**: When you type a URL, your browser first contacts a DNS server to translate the domain name into an IP address.

2. **TCP Connection**: The browser initiates a TCP connection with the server using a three-way handshake.

3. **HTTP Request**: Once the connection is established, the browser sends an HTTP request to the server.

4. **Server Processing**: The server processes the request and prepares a response.

5. **HTTP Response**: The server sends back an HTTP response, which typically includes a status code and the requested content.

6. **Rendering**: The browser renders the received content (usually HTML, CSS, and JavaScript).

### What We'll Build

Throughout this series, we'll build:

1. A basic TCP server
2. An HTTP server capable of handling GET and POST requests
3. A simple routing system
4. Support for serving static files
5. Basic HTTPS functionality

### Series Outline

1. Part 0: Foundations and Prerequisites (this part)
2. Part 1: Building a Basic TCP Server
3. Part 2: HTTP Basics - Parsing Requests
4. Part 3: HTTP Responses and Status Codes
5. Part 4: Implementing GET and POST Methods
6. Part 5: Simple Routing System
7. Part 6: Serving Static Files
8. Part 7: Introduction to HTTPS
9. Part 8: Implementing HTTPS in Our Server

### Conclusion

We've covered the basic concepts and prerequisites for our journey to build HTTP from scratch. In the next part, we'll start by implementing a basic TCP server, which will serve as the foundation for our HTTP server.

Remember, building a production-ready web server is complex and requires careful consideration of security, performance, and edge cases. This series is meant for educational purposes to understand the underlying principles of HTTP and network programming.

Stay tuned for Part 1, where we'll dive into creating our first TCP server!