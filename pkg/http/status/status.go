package status

const (
	OK                   = 200
	Created              = 201
	Accepted             = 202
	NoContent            = 204
	MovedPermanently     = 301
	Found                = 302
	BadRequest           = 400
	Unauthorized         = 401
	Forbidden            = 403
	NotFound             = 404
	MethodNotAllowed     = 405
	StatusRequestTimeout = 408
	IamATeaPot           = 418
	InternalServerError  = 500
	NotImplemented       = 501
	BadGateway           = 502
	ServiceUnavailable   = 503
)

var statusText = map[int]string{
	OK:                   "OK",
	Created:              "Created",
	Accepted:             "Accepted",
	NoContent:            "No Content",
	MovedPermanently:     "Moved Permanently",
	Found:                "Found",
	BadRequest:           "Bad Request",
	Unauthorized:         "Unauthorized",
	Forbidden:            "Forbidden",
	NotFound:             "Not Found",
	MethodNotAllowed:     "Method Not Allowed",
	IamATeaPot:           "I'm a teapot",
	InternalServerError:  "Internal Server Error",
	NotImplemented:       "Not Implemented",
	BadGateway:           "Bad Gateway",
	ServiceUnavailable:   "Service Unavailable",
	StatusRequestTimeout: "Request Timeout",
}

func Text(code int) string {
	return statusText[code]
}
