package http

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/appyzdl/Netrunner/pkg/http"
)

func TestHandleEcho(t *testing.T) {
	// Create a test request
	testBody := []byte("Hello, Netrunner!")
	req := &http.Request{
		Method:  "POST",
		Path:    "/echo",
		Version: "HTTP/1.1",
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "17",
		},
		Body: testBody,
	}

	// Call the handler
	resp := handleEcho(req)

	// Check the response status
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
	if resp.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got '%s'", resp.StatusText)
	}

	// Check the response headers
	contentType := resp.Headers["Content-Type"]
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", contentType)
	}
	contentLength := resp.Headers["Content-Length"]
	if contentLength != "17" {
		t.Errorf("Expected Content-Length '17', got '%s'", contentLength)
	}

	// Check the response body
	if !bytes.Equal(resp.Body, testBody) {
		t.Errorf("Expected body '%s', got '%s'", testBody, resp.Body)
	}
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
