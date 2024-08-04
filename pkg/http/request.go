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
			break
		}
		parts := strings.SplitN(lines[i], ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", lines[i])
		}
		request.Headers[parts[0]] = parts[1]
	}

	return request, nil
}
