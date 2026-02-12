// Package nwfetch provides a high-level fetch client for the NWEP/WEB/1
// protocol, built on top of github.com/usenwep/nwep-go. It manages connection
// pooling, identity, and request construction so that callers can interact with
// WEB/1 servers using a familiar HTTP-style API while keeping the underlying
// nwep types (nwep.Keypair, nwep.Header, nwep.Notification, etc.) fully
// transparent - nothing is re-wrapped.
//
// Quick start using the default client:
//
//	nwfetch.Init()
//	defer nwfetch.Close()
//	resp, err := nwfetch.Get("web://addr/hello")
//	fmt.Println(resp.String())
//
// Or using the fluent request builder:
//
//	resp, err := nwfetch.New("web://addr/path").
//	    Method("write").
//	    Header("content-type", "application/json").
//	    Body(jsonData).
//	    Do()
//
// For long-lived programs or programs that need a stable identity, create a
// Client with explicit options:
//
//	client, err := nwfetch.NewClient(
//	    nwfetch.WithKeypair(kp),
//	    nwfetch.WithTimeout(10 * time.Second),
//	)
//	defer client.Close()
//	resp, err := client.Get("web://addr/hello")
//
// nwfetch has zero external dependencies beyond nwep-go.
package nwfetch

import (
	"sync"

	"github.com/usenwep/nwep-go"
)

var (
	defaultClient *Client
	defaultMu     sync.Mutex
)

// Init initializes the NWEP runtime and creates the default client with an
// ephemeral Ed25519 identity. It must be called before using the top-level
// convenience functions (Get, Post, Do) or the request builder's Do method.
//
// This function returns a non-nil error if NWEP runtime initialization fails
// or if keypair generation fails. Calling Init more than once replaces the
// previous default client.
func Init() error {
	if err := nwep.Init(); err != nil {
		return err
	}

	c, err := NewClient()
	if err != nil {
		return err
	}

	defaultMu.Lock()
	defaultClient = c
	defaultMu.Unlock()
	return nil
}

// Version returns the underlying NWEP library version string.
func Version() string {
	return nwep.Version()
}

// Default returns the default client created by Init. It panics if Init has
// not been called.
func Default() *Client {
	defaultMu.Lock()
	c := defaultClient
	defaultMu.Unlock()
	if c == nil {
		panic("nwfetch: Init() must be called before using the default client")
	}
	return c
}

// Close closes the default client, releasing all pooled connections and
// clearing the ephemeral keypair. After Close returns, the top-level
// convenience functions will panic until Init is called again.
func Close() {
	defaultMu.Lock()
	c := defaultClient
	defaultClient = nil
	defaultMu.Unlock()
	if c != nil {
		c.Close()
	}
}

// Get performs a "read" request to the given URL using the default client.
// It is shorthand for Default().Get(url).
func Get(url string) (*Response, error) {
	return Default().Get(url)
}

// Post performs a "write" request to the given URL with the provided body
// using the default client. It is shorthand for Default().Post(url, body).
func Post(url string, body []byte) (*Response, error) {
	return Default().Post(url, body)
}

// Do executes a Request using the default client. It is shorthand for
// Default().Do(req).
func Do(req *Request) (*Response, error) {
	return Default().Do(req)
}
