package http_test

import (
	"testing"

	"github.com/appyzdl/Netrunner/pkg/http"
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
