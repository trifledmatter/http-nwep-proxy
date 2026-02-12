package nwfetch

import (
	"time"

	"github.com/usenwep/nwep-go"
)

// Request is a fluent builder for constructing fetch requests. A Request is
// created with New, configured with chained method calls, and executed with
// Do (default client) or DoWith (explicit client).
//
// The zero-value method is "read". All builder methods return the same
// *Request so calls can be chained:
//
//	resp, err := nwfetch.New("web://addr/path").
//	    Method("write").
//	    Header("content-type", "application/json").
//	    Body(jsonData).
//	    Do()
//
// Request instances are not safe for concurrent use and must not be reused
// after calling Do or DoWith.
type Request struct {
	url     string
	method  string
	headers []nwep.Header
	body    []byte
	timeout time.Duration
}

// New creates a new Request for the given WEB/1 URL with the default method
// "read". The URL must be in web:// format (e.g. "web://addr/path").
func New(url string) *Request {
	return &Request{
		url:    url,
		method: nwep.MethodRead,
	}
}

// Method sets the request method. Valid methods are "read", "write", "update",
// and "delete". See the nwep.Method* constants for the full set.
func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

// Header appends a header to the request. Multiple headers with the same name
// are allowed.
func (r *Request) Header(name, value string) *Request {
	r.headers = append(r.headers, nwep.Header{Name: name, Value: value})
	return r
}

// Body sets the request body. Passing nil sends a request with no body.
func (r *Request) Body(body []byte) *Request {
	r.body = body
	return r
}

// Timeout sets a per-request timeout that overrides the client's default. A
// zero value means no per-request timeout.
func (r *Request) Timeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

// Do executes the request using the default client. Init must have been called
// before using this method.
func (r *Request) Do() (*Response, error) {
	return Default().Do(r)
}

// DoWith executes the request using the given client.
func (r *Request) DoWith(c *Client) (*Response, error) {
	return c.Do(r)
}
