package http

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
	TLS     *tls.ConnectionState
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}

func ParseRequest(data []byte, tlsConn *tls.ConnectionState) (*Request, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	// Read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %v", err)
	}
	requestLine = strings.TrimSpace(requestLine)

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	request := &Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
		Headers: make(map[string]string),
		TLS:     tlsConn,
	}

	// Read headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		request.Headers[key] = value
	}

	// Read body if present
	contentLength := request.Headers["Content-Length"]
	if contentLength != "" {
		// Implementation for reading body based on Content-Length
		// This is a simplified version and may need to be enhanced
		bodyBuffer := make([]byte, len(data))
		n, err := reader.Read(bodyBuffer)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
		request.Body = bodyBuffer[:n]
	}

	return request, nil
}
