package nwfetch

import nwep "github.com/usenwep/nwep-go"

// WEB/1 response status constants, re-exported from nwep for convenience so
// that callers do not need to import nwep directly for status comparisons.
// Success statuses can be checked with nwep.StatusIsSuccess; error statuses
// with nwep.StatusIsError.
const (
	// StatusOK indicates the request was processed successfully.
	StatusOK = nwep.StatusOK

	// StatusCreated indicates a new resource was created as a result of the
	// request.
	StatusCreated = nwep.StatusCreated

	// StatusAccepted indicates the request has been accepted for processing,
	// but processing has not completed.
	StatusAccepted = nwep.StatusAccepted

	// StatusNoContent indicates the request succeeded but there is no
	// response body.
	StatusNoContent = nwep.StatusNoContent

	// StatusBadRequest indicates the request was malformed or contained
	// invalid parameters.
	StatusBadRequest = nwep.StatusBadRequest

	// StatusUnauthorized indicates the peer has not provided valid
	// authentication credentials.
	StatusUnauthorized = nwep.StatusUnauthorized

	// StatusForbidden indicates the peer is authenticated but lacks
	// permission for the requested operation.
	StatusForbidden = nwep.StatusForbidden

	// StatusNotFound indicates the requested path does not exist on the
	// server.
	StatusNotFound = nwep.StatusNotFound

	// StatusConflict indicates the request conflicts with the current state
	// of the resource (e.g. a duplicate write).
	StatusConflict = nwep.StatusConflict

	// StatusRateLimited indicates the peer has exceeded the allowed request
	// rate. The peer should back off before retrying.
	StatusRateLimited = nwep.StatusRateLimited

	// StatusInternalError indicates an unexpected server-side failure.
	StatusInternalError = nwep.StatusInternalError

	// StatusUnavailable indicates the server is temporarily unable to handle
	// the request (e.g. during startup or overload).
	StatusUnavailable = nwep.StatusUnavailable
)
