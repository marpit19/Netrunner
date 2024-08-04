package http

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

func NotFoundResponse() *Response {
	resp := NewResponse()
	resp.StatusCode = 404
	resp.StatusText = "Not Found"
	resp.SetBody([]byte("404 - Not Found"))
	return resp
}
