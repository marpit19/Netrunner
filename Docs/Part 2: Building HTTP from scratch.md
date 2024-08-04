# Netrunner: Building HTTP from the Ground Up
## Part 2: HTTP Basics - Parsing Requests

Welcome to Part 2 of our Netrunner series! In this part, we'll build on our TCP server to start handling HTTP requests. We'll focus on understanding the structure of HTTP requests and implementing a parser for them.

### HTTP Request Structure

Before we start coding, let's review the structure of an HTTP request:

1. Request Line: Contains the HTTP method, request target, and HTTP version
2. Headers: Key-value pairs providing additional information about the request
3. Empty line: Separates headers from the body
4. Body (optional): Contains the payload of the request (e.g., form data for POST requests)

A typical HTTP request looks like this:

```
GET /index.html HTTP/1.1
Host: www.example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8

```

### Setting Up Our Project

Let's create new files for our HTTP implementation:

```bash
mkdir -p netrunner/pkg/http
touch netrunner/pkg/http/request.go
mkdir -p netrunner/test/http
touch netrunner/test/http/request_test.go
```

### Defining the Request Structure

In `netrunner/pkg/http/request.go`, let's define our HTTP request structure:

```go
package http

import (
	"fmt"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}
```

### Writing Our First Test

In `netrunner/test/http/request_test.go`, let's write a test for parsing a simple HTTP request:

```go
package http_test

import (
	"testing"

	"github.com/yourusername/netrunner/pkg/http"
)

func TestParseRequest(t *testing.T) {
	rawRequest := "GET /index.html HTTP/1.1\r\n" +
		"Host: www.example.com\r\n" +
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36\r\n" +
		"\r\n"

	request, err := http.ParseRequest([]byte(rawRequest))
	if err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}

	if request.Method != "GET" {
		t.Errorf("Expected method GET, got %s", request.Method)
	}

	if request.Path != "/index.html" {
		t.Errorf("Expected path /index.html, got %s", request.Path)
	}

	if request.Version != "HTTP/1.1" {
		t.Errorf("Expected version HTTP/1.1, got %s", request.Version)
	}

	if len(request.Headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(request.Headers))
	}

	if request.Headers["Host"] != "www.example.com" {
		t.Errorf("Expected Host header www.example.com, got %s", request.Headers["Host"])
	}

	expectedUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	if request.Headers["User-Agent"] != expectedUserAgent {
		t.Errorf("Expected User-Agent header %s, got %s", expectedUserAgent, request.Headers["User-Agent"])
	}
}
```

### Implementing the Request Parser

Now, let's implement the `ParseRequest` function in `request.go`:

```go
func ParseRequest(data []byte) (*Request, error) {
	request := NewRequest()

	lines := strings.Split(string(data), "\r\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid request: too few lines")
	}

	// Parse request line
	parts := strings.Split(lines[0], " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", lines[0])
	}
	request.Method = parts[0]
	request.Path = parts[1]
	request.Version = parts[2]

	// Parse headers
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			break // Empty line signifies end of headers
		}
		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", lines[i])
		}
		request.Headers[parts[0]] = parts[1]
	}

	return request, nil
}
```

This implementation does the following:
1. Splits the raw request into lines
2. Parses the first line to extract the method, path, and HTTP version
3. Parses subsequent lines as headers until it encounters an empty line

### Running the Test

To run the test, navigate to your project root and run:

```bash
go test ./test/http
```

If everything is implemented correctly, the test should pass.

### Handling Request Body

To handle request bodies (important for POST requests), let's update our `ParseRequest` function:

```go
func ParseRequest(data []byte) (*Request, error) {
	request := NewRequest()

	parts := bytes.SplitN(data, []byte("\r\n\r\n"), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid request: no body separator found")
	}

	lines := strings.Split(string(parts[0]), "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid request: empty request")
	}

	// Parse request line
	requestLineParts := strings.Split(lines[0], " ")
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", lines[0])
	}
	request.Method = requestLineParts[0]
	request.Path = requestLineParts[1]
	request.Version = requestLineParts[2]

	// Parse headers
	for i := 1; i < len(lines); i++ {
		headerParts := strings.SplitN(lines[i], ": ", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", lines[i])
		}
		request.Headers[headerParts[0]] = headerParts[1]
	}

	// Set body
	request.Body = parts[1]

	return request, nil
}
```

### Testing Request Body Parsing

Add a new test to `request_test.go`:

```go
func TestParseRequestWithBody(t *testing.T) {
	rawRequest := "POST /submit HTTP/1.1\r\n" +
		"Host: www.example.com\r\n" +
		"Content-Type: application/x-www-form-urlencoded\r\n" +
		"Content-Length: 27\r\n" +
		"\r\n" +
		"username=johndoe&password=123"

	request, err := http.ParseRequest([]byte(rawRequest))
	if err != nil {
		t.Fatalf("Failed to parse request: %v", err)
	}

	if request.Method != "POST" {
		t.Errorf("Expected method POST, got %s", request.Method)
	}

	if request.Path != "/submit" {
		t.Errorf("Expected path /submit, got %s", request.Path)
	}

	expectedBody := "username=johndoe&password=123"
	if string(request.Body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(request.Body))
	}
}
```

### Conclusion

We've now implemented a basic HTTP request parser. Key takeaways from this part:

1. We defined a structure to represent HTTP requests
2. We implemented parsing for the request line, headers, and body
3. We wrote tests to verify our parser's functionality

In the next part, we'll focus on generating HTTP responses and implementing status codes. We'll also start to integrate this with our TCP server from Part 1 to create a functional (albeit basic) HTTP server.

Remember, this is a simplified implementation for educational purposes. Production-grade HTTP servers need to handle various edge cases, support all HTTP methods, and implement more robust error handling.

Stay tuned for Part 3, where we'll dive into HTTP responses and status codes!