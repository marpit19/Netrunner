package http

import (
	"fmt"
	"time"
)

// LoggingMiddleware logs information about each request
func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(req *Request) *Response {
		start := time.Now()
		resp := next(req)
		duration := time.Since(start)
		fmt.Printf("%s %s - %d (%v)\n", req.Method, req.Path, resp.StatusCode, duration)
		return resp
	}
}
