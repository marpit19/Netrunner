package http

import (
	"fmt"
	"strings"

	"github.com/appyzdl/Netrunner/pkg/http/status"
)

type (
	HandlerFunc    func(*Request) *Response
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

type Router struct {
	routes     map[string]map[string]HandlerFunc
	middleware []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{
		routes:     make(map[string]map[string]HandlerFunc),
		middleware: []MiddlewareFunc{},
	}
}

func (r *Router) Use(mw MiddlewareFunc) {
	r.middleware = append(r.middleware, mw)
}

func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) HandleRequest(req *Request) *Response {
	if req.TLS == nil && r.shouldRedirectToHTTPS(req) {
		return r.redirectToHTTPS(req)
	}

	if handlers, ok := r.routes[req.Method]; ok {
		if handler, ok := handlers[req.Path]; ok {
			// middleware
			for i := len(r.middleware) - 1; i >= 0; i-- {
				handler = r.middleware[i](handler)
			}
			return handler(req)
		}
	}
	return NotFoundResponse()
}

func (r *Router) shouldRedirectToHTTPS(req *Request) bool {
	return !strings.HasPrefix(req.Path, "/static/")
}

func (r *Router) redirectToHTTPS(req *Request) *Response {
	resp := NewResponse()
	resp.StatusCode = status.MovedPermanently
	resp.StatusText = StatusText(status.MovedPermanently)
	httpsURL := fmt.Sprintf("https://%s%s", req.Headers["Host"], req.Path)
	resp.SetHeader("Location", httpsURL)
	return resp
}

func NotFoundResponse() *Response {
	resp := NewResponse()
	resp.StatusCode = 404
	resp.StatusText = "Not Found"
	resp.SetBody([]byte("404 - Not Found"))
	return resp
}
