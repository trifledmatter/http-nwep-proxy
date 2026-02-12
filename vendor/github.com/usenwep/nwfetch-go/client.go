package nwfetch

import (
	"time"

	"github.com/usenwep/nwep-go"
)

// Client is a high-level fetch client for the NWEP/WEB/1 protocol. It owns an
// Ed25519 keypair (auto-generated or provided), a connection pool, and
// per-client settings such as timeout and notification callback.
//
// A Client is created with NewClient, configured with ClientOption functions,
// and must be closed with Close when no longer needed. The Do method is the
// primary entry point for executing requests; Get and Post are convenience
// wrappers around Do.
//
// All methods on Client are safe for concurrent use. Connections are pooled and
// reused automatically - NWEP multiplexes streams over a single connection per
// server, so no pool sizing is required.
type Client struct {
	keypair  *nwep.Keypair
	ownsKey  bool
	pool     *connPool
	timeout  time.Duration
	settings *nwep.Settings
}

// ClientOption configures a Client. Options are applied in order during
// NewClient. See the With* functions for available options.
type ClientOption func(*clientConfig)

type clientConfig struct {
	keypair  *nwep.Keypair
	seed     *[32]byte
	timeout  time.Duration
	settings *nwep.Settings
	onNotify func(*nwep.Notification)
	poolSize int
}

// WithKeypair sets the client identity to an existing Ed25519 keypair. The
// caller retains ownership of the keypair - Close will not clear it. If both
// WithKeypair and WithSeed are provided, the last one applied wins.
func WithKeypair(kp *nwep.Keypair) ClientOption {
	return func(c *clientConfig) { c.keypair = kp }
}

// WithSeed derives a deterministic Ed25519 keypair from the given 32-byte
// seed. The Client takes ownership of the derived keypair and will clear it
// on Close.
func WithSeed(seed [32]byte) ClientOption {
	return func(c *clientConfig) { c.seed = &seed }
}

// WithTimeout sets the default timeout applied to all requests made through
// this client. A zero value means no timeout (the NWEP default applies).
func WithTimeout(d time.Duration) ClientOption {
	return func(c *clientConfig) { c.timeout = d }
}

// WithSettings sets the NWEP protocol settings (max streams, max message size,
// compression, etc.) used when establishing new connections. Zero-valued fields
// in the Settings struct are ignored by the underlying library.
func WithSettings(s nwep.Settings) ClientOption {
	return func(c *clientConfig) { c.settings = &s }
}

// WithOnNotify registers a callback that is invoked when the server sends an
// unsolicited notification. The callback is called on the nwep event loop
// goroutine, so it must not block.
func WithOnNotify(fn func(*nwep.Notification)) ClientOption {
	return func(c *clientConfig) { c.onNotify = fn }
}

// WithPoolSize is reserved for future use. NWEP multiplexes streams over a
// single connection per server, so only one connection per address is
// maintained regardless of this setting.
func WithPoolSize(n int) ClientOption {
	return func(c *clientConfig) { c.poolSize = n }
}

// NewClient creates a new Client with the given options. Options are applied
// in order.
//
// If no identity option is provided (WithKeypair or WithSeed), a random
// ephemeral Ed25519 keypair is generated. This function returns a non-nil
// error if keypair generation or derivation fails.
func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{}
	for _, o := range opts {
		o(cfg)
	}

	var kp *nwep.Keypair
	var ownsKey bool

	switch {
	case cfg.keypair != nil:
		kp = cfg.keypair
	case cfg.seed != nil:
		var err error
		kp, err = nwep.KeypairFromSeed(*cfg.seed)
		if err != nil {
			return nil, err
		}
		ownsKey = true
	default:
		var err error
		kp, err = nwep.GenerateKeypair()
		if err != nil {
			return nil, err
		}
		ownsKey = true
	}

	var nwepOpts []nwep.ClientOption
	if cfg.onNotify != nil {
		nwepOpts = append(nwepOpts, nwep.WithOnNotify(cfg.onNotify))
	}
	if cfg.settings != nil {
		nwepOpts = append(nwepOpts, nwep.WithClientSettings(*cfg.settings))
	}

	c := &Client{
		keypair:  kp,
		ownsKey:  ownsKey,
		pool:     newConnPool(kp, nwepOpts),
		timeout:  cfg.timeout,
		settings: cfg.settings,
	}
	return c, nil
}

// Do executes a Request and returns the Response. It parses the request URL,
// obtains a pooled connection to the target server (connecting if necessary),
// and performs the fetch.
//
// On transport errors the failing connection is removed from the pool and the
// error is wrapped in an *Error with the Op field set to "parse", "connect",
// or "fetch" depending on where the failure occurred. Protocol-level errors
// (e.g. "not_found") are not treated as transport errors - they are returned
// as a successful *Response whose status can be inspected with IsSuccess,
// IsError, or StatusError.
func (c *Client) Do(req *Request) (*Response, error) {
	parsed, err := nwep.URLParse(NormalizeURL(req.url))
	if err != nil {
		return nil, &Error{Op: "parse", URL: req.url, Err: err}
	}

	conn, key, err := c.pool.get(parsed)
	if err != nil {
		return nil, &Error{Op: "connect", URL: req.url, Err: err}
	}

	nr, err := conn.FetchWithHeaders(req.method, parsed.Path, req.body, req.headers)
	if err != nil {
		c.pool.remove(key)
		return nil, &Error{Op: "fetch", URL: req.url, Err: err}
	}

	return responseFromNWEP(nr), nil
}

// Get performs a "read" request to the given URL. It is shorthand for
// Do(New(url)).
func (c *Client) Get(url string) (*Response, error) {
	return c.Do(New(url))
}

// Post performs a "write" request to the given URL with the provided body. It
// is shorthand for Do(New(url).Method("write").Body(body)).
func (c *Client) Post(url string, body []byte) (*Response, error) {
	return c.Do(New(url).Method(nwep.MethodWrite).Body(body))
}

// Close closes all pooled connections and, if the Client owns its keypair
// (generated or derived from seed), securely clears the key material. After
// Close returns the Client must not be used.
func (c *Client) Close() {
	c.pool.closeAll()
	if c.ownsKey && c.keypair != nil {
		c.keypair.Clear()
	}
}
