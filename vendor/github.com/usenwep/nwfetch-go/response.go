package nwfetch

import (
	"strconv"
	"time"

	"github.com/usenwep/nwep-go"
)

// Response wraps a WEB/1 protocol response. The Status field contains the
// server's status string (e.g. "ok", "not_found"), StatusDetails carries an
// optional human-readable elaboration, Headers holds the response headers, and
// Body holds the response payload.
//
// Use IsOK, IsSuccess, or IsError for quick status checks. Use StatusError to
// convert an error response into a Go error suitable for returning up a call
// chain.
type Response struct {
	// Status is the WEB/1 status string (e.g. "ok", "created",
	// "not_found"). See the nwep.Status* constants for defined values.
	Status string

	// StatusDetails is an optional server-provided elaboration on the
	// status, typically a human-readable error message.
	StatusDetails string

	// Headers contains the response headers. Use the Header method for
	// convenient lookup by name.
	Headers []nwep.Header

	// Body is the raw response payload. It may be nil for responses with
	// no content.
	Body []byte
}

// IsOK reports whether the status is exactly "ok".
func (r *Response) IsOK() bool {
	return r.Status == nwep.StatusOK
}

// IsSuccess reports whether the status is any success status ("ok", "created",
// "accepted", "no_content").
func (r *Response) IsSuccess() bool {
	return nwep.StatusIsSuccess(r.Status)
}

// IsError reports whether the status is an error status ("bad_request",
// "unauthorized", "forbidden", "not_found", etc.).
func (r *Response) IsError() bool {
	return nwep.StatusIsError(r.Status)
}

// Header returns the value of the first header matching name, and a boolean
// indicating whether the header was found. Header names are case-sensitive.
func (r *Response) Header(name string) (string, bool) {
	for _, h := range r.Headers {
		if h.Name == name {
			return h.Value, true
		}
	}
	return "", false
}

// String returns the response body as a string.
func (r *Response) String() string {
	return string(r.Body)
}

// RetryAfter returns the duration a client should wait before retrying, as
// indicated by the "retry-after" header. The WEB/1 spec requires this header
// on "rate_limited" responses. If the header is missing or cannot be parsed
// as an integer number of seconds, ok is false.
func (r *Response) RetryAfter() (d time.Duration, ok bool) {
	v, found := r.Header("retry-after")
	if !found {
		return 0, false
	}
	secs, err := strconv.Atoi(v)
	if err != nil || secs < 0 {
		return 0, false
	}
	return time.Duration(secs) * time.Second, true
}

// StatusError returns nil if the response has a success status. For error
// statuses it returns a *StatusError containing the status, details, and body.
// This is useful for propagating server errors up a call chain:
//
//	resp, err := client.Get("web://addr/resource")
//	if err != nil { return err }
//	if err := resp.StatusError(); err != nil { return err }
func (r *Response) StatusError() error {
	if r.IsSuccess() {
		return nil
	}
	return &StatusError{
		Status:        r.Status,
		StatusDetails: r.StatusDetails,
		Body:          r.Body,
	}
}

func responseFromNWEP(nr *nwep.Response) *Response {
	return &Response{
		Status:        nr.Status,
		StatusDetails: nr.StatusDetails,
		Headers:       nr.Headers,
		Body:          nr.Body,
	}
}
