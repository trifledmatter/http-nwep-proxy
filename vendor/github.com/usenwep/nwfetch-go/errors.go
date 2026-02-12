package nwfetch

import (
	"errors"
	"fmt"
)

// Error represents a transport-level error that occurred while executing a
// fetch operation. It wraps the underlying error with context about which
// operation failed and which URL was being accessed.
//
// The Op field indicates the stage of the request that failed: "parse" for URL
// parsing, "connect" for connection establishment, or "fetch" for the actual
// data exchange. The wrapped Err can be inspected with errors.Is and
// errors.As for more detail.
type Error struct {
	// Op is the operation that failed (e.g. "parse", "connect", "fetch").
	Op string

	// URL is the WEB/1 URL that was being accessed when the error occurred.
	URL string

	// Err is the underlying error from the nwep library.
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("nwfetch: %s %s: %s", e.Op, e.URL, e.Err)
}

func (e *Error) Unwrap() error { return e.Err }

// StatusError represents a protocol-level error response from a WEB/1 server.
// Unlike Error, which indicates a transport failure, a StatusError means the
// server was reached and responded with a non-success status.
//
// StatusError values are returned by Response.StatusError and can be matched
// with the Is* helper functions (IsNotFound, IsUnauthorized, etc.) or
// inspected directly via errors.As:
//
//	var se *nwfetch.StatusError
//	if errors.As(err, &se) {
//	    fmt.Println(se.Status, se.StatusDetails)
//	}
type StatusError struct {
	// Status is the WEB/1 status string (e.g. "not_found", "forbidden").
	Status string

	// StatusDetails is an optional server-provided elaboration, typically
	// a human-readable error message.
	StatusDetails string

	// Body is the raw response body, which may contain additional error
	// information from the server.
	Body []byte
}

func (e *StatusError) Error() string {
	if e.StatusDetails != "" {
		return fmt.Sprintf("nwfetch: server returned %s: %s", e.Status, e.StatusDetails)
	}
	return fmt.Sprintf("nwfetch: server returned %s", e.Status)
}

// IsBadRequest reports whether err is or wraps a *StatusError with status
// "bad_request".
func IsBadRequest(err error) bool { return hasStatus(err, "bad_request") }

// IsUnauthorized reports whether err is or wraps a *StatusError with status
// "unauthorized".
func IsUnauthorized(err error) bool { return hasStatus(err, "unauthorized") }

// IsForbidden reports whether err is or wraps a *StatusError with status
// "forbidden".
func IsForbidden(err error) bool { return hasStatus(err, "forbidden") }

// IsNotFound reports whether err is or wraps a *StatusError with status
// "not_found".
func IsNotFound(err error) bool { return hasStatus(err, "not_found") }

// IsConflict reports whether err is or wraps a *StatusError with status
// "conflict".
func IsConflict(err error) bool { return hasStatus(err, "conflict") }

// IsRateLimited reports whether err is or wraps a *StatusError with status
// "rate_limited".
func IsRateLimited(err error) bool { return hasStatus(err, "rate_limited") }

// IsInternalError reports whether err is or wraps a *StatusError with status
// "internal_error".
func IsInternalError(err error) bool { return hasStatus(err, "internal_error") }

// IsUnavailable reports whether err is or wraps a *StatusError with status
// "unavailable".
func IsUnavailable(err error) bool { return hasStatus(err, "unavailable") }

func hasStatus(err error, status string) bool {
	var se *StatusError
	return errors.As(err, &se) && se.Status == status
}
