package http

import (
	"bytes"
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
