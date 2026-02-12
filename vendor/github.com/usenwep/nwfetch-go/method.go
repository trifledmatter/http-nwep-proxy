package nwfetch

import nwep "github.com/usenwep/nwep-go"

// WEB/1 request method constants, re-exported from nwep for convenience so
// that callers do not need to import nwep directly. These are the method
// strings used with Request.Method and the Client convenience methods.
const (
	// MethodRead requests retrieval of a resource. Idempotent and safe for
	// 0-RTT. This is the default method used by New and Client.Get.
	MethodRead = nwep.MethodRead

	// MethodWrite requests creation of a new resource. Not idempotent.
	// This is the method used by Client.Post.
	MethodWrite = nwep.MethodWrite

	// MethodUpdate requests modification of an existing resource. Not
	// idempotent.
	MethodUpdate = nwep.MethodUpdate

	// MethodDelete requests removal of a resource. Idempotent.
	MethodDelete = nwep.MethodDelete
)
