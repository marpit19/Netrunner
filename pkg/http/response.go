package http

import (
	"fmt"
	"strings"

	"github.com/appyzdl/Netrunner/pkg/http/status"
)

type Response struct {
	Version    string
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func NewResponse() *Response {
	return &Response{
		Version: "HTTP/1.1",
		Headers: make(map[string]string),
	}
}

func (r *Response) SetStatus(code int) {
	r.StatusCode = code
}

func (r *Response) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *Response) SetBody(body []byte) {
	r.Body = body
	r.SetHeader("Content-Length", fmt.Sprintf("%d", len(body)))
}

func (r *Response) Write() []byte {
	var builder strings.Builder

	statusText := status.Text(r.StatusCode)
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Version, r.StatusCode, statusText))

	for key, value := range r.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	builder.WriteString("\r\n")
	return append([]byte(builder.String()), r.Body...)
}
