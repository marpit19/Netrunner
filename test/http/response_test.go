package http_test

import (
	"strings"
	"testing"

	"github.com/appyzdl/Netrunner/pkg/http"
	"github.com/appyzdl/Netrunner/pkg/http/status"
)

func TestResponse(t *testing.T) {
	response := http.NewResponse()
	response.SetStatus(status.OK)
	response.SetHeader("Content-Type", "text/plain")
	response.SetBody([]byte("Hello, Netrunner!"))

	rawResponse := response.Write()
	lines := strings.Split(string(rawResponse), "\r\n")

	// Check status line
	if lines[0] != "HTTP/1.1 200 OK" {
		t.Errorf("Expected status line 'HTTP/1.1 200 OK', got '%s'", lines[0])
	}

	// Check headers
	expectedHeaders := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": "17",
	}

	for i := 1; i < len(lines)-2; i++ {
		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) != 2 {
			t.Errorf("Invalid header line: %s", lines[i])
			continue
		}
		if expectedHeaders[parts[0]] != parts[1] {
			t.Errorf("Expected header '%s: %s', got '%s'", parts[0], expectedHeaders[parts[0]], lines[i])
		}
	}

	// Check body
	body := lines[len(lines)-1]
	if body != "Hello, Netrunner!" {
		t.Errorf("Expected body 'Hello, Netrunner!', got '%s'", body)
	}
}
