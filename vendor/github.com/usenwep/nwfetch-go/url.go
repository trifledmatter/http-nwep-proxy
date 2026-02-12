package nwfetch

import "strings"

// DefaultPort is the default port for the WEB/1 protocol.
const DefaultPort = "6937"

// NormalizeURL takes a user-provided URL and ensures it has the web:// scheme,
// brackets around the address, and a port. Supported input forms:
//
//	web://[addr]:port/path  → unchanged
//	web://addr:port/path    → web://[addr]:port/path
//	web://[addr]/path       → web://[addr]:6937/path
//	web://addr/path         → web://[addr]:6937/path
//	web://addr              → web://[addr]:6937/
func NormalizeURL(raw string) string {
	// Strip scheme.
	rest := raw
	if strings.HasPrefix(rest, "web://") {
		rest = rest[len("web://"):]
	}

	if len(rest) == 0 {
		return raw
	}

	// Split host part from path.
	var host, path string
	if rest[0] == '[' {
		// Bracketed form: web://[addr]:port/path or web://[addr]/path
		end := strings.Index(rest, "]")
		if end == -1 {
			// Malformed, just return as-is and let nwep report the error.
			return raw
		}
		addr := rest[1:end]          // inside brackets
		afterBracket := rest[end+1:] // everything after ']'

		if len(afterBracket) == 0 {
			// web://[addr]
			host = "[" + addr + "]:" + DefaultPort
			path = "/"
		} else if afterBracket[0] == ':' {
			// web://[addr]:port... — already has port
			slashIdx := strings.Index(afterBracket, "/")
			if slashIdx == -1 {
				host = "[" + addr + "]" + afterBracket
				path = "/"
			} else {
				host = "[" + addr + "]" + afterBracket[:slashIdx]
				path = afterBracket[slashIdx:]
			}
		} else if afterBracket[0] == '/' {
			// web://[addr]/path — no port
			host = "[" + addr + "]:" + DefaultPort
			path = afterBracket
		} else {
			return raw
		}
	} else {
		// Unbracketed form: web://addr:port/path or web://addr/path
		slashIdx := strings.Index(rest, "/")
		var hostPart string
		if slashIdx == -1 {
			hostPart = rest
			path = "/"
		} else {
			hostPart = rest[:slashIdx]
			path = rest[slashIdx:]
		}

		if strings.Contains(hostPart, ":") {
			// Has port already: addr:port → [addr]:port
			parts := strings.SplitN(hostPart, ":", 2)
			host = "[" + parts[0] + "]:" + parts[1]
		} else {
			// No port: addr → [addr]:6937
			host = "[" + hostPart + "]:" + DefaultPort
		}
	}

	return "web://" + host + path
}
