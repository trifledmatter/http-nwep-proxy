package nwfetch

import (
	"sync"

	"github.com/usenwep/nwep-go"
)

// connPool manages a set of reusable nwep.Client connections keyed by the
// base58-encoded server address (IP + port + nodeID). Because NWEP multiplexes
// all streams over a single connection, only one connection per server is
// maintained.
//
// All methods are safe for concurrent use.
type connPool struct {
	mu      sync.Mutex
	conns   map[string]*nwep.Client
	keypair *nwep.Keypair
	opts    []nwep.ClientOption
}

func newConnPool(kp *nwep.Keypair, opts []nwep.ClientOption) *connPool {
	return &connPool{
		conns:   make(map[string]*nwep.Client),
		keypair: kp,
		opts:    opts,
	}
}

// get returns an existing connection for the URL's server address, or dials
// and connects a new one. The returned key can be passed to remove if the
// connection later fails.
//
// When two goroutines race to connect to the same server, the second one
// detects the duplicate during the final lock acquisition and discards its
// connection in favor of the one already stored.
func (p *connPool) get(u *nwep.URL) (*nwep.Client, string, error) {
	key, err := nwep.AddrEncode(&u.Addr)
	if err != nil {
		return nil, "", err
	}

	p.mu.Lock()
	if c, ok := p.conns[key]; ok {
		p.mu.Unlock()
		return c, key, nil
	}
	p.mu.Unlock()

	c, err := nwep.NewClient(p.keypair, p.opts...)
	if err != nil {
		return nil, "", err
	}

	rawURL, err := nwep.URLFormat(u)
	if err != nil {
		c.Close()
		return nil, "", err
	}

	if err := c.Connect(rawURL); err != nil {
		c.Close()
		return nil, "", err
	}

	p.mu.Lock()
	// Check again in case another goroutine connected concurrently.
	if existing, ok := p.conns[key]; ok {
		p.mu.Unlock()
		c.Close()
		return existing, key, nil
	}
	p.conns[key] = c
	p.mu.Unlock()

	return c, key, nil
}

// remove removes a connection by key and closes it. This is called when a
// fetch fails so that subsequent requests will establish a fresh connection.
func (p *connPool) remove(key string) {
	p.mu.Lock()
	if c, ok := p.conns[key]; ok {
		delete(p.conns, key)
		p.mu.Unlock()
		c.Close()
		return
	}
	p.mu.Unlock()
}

// closeAll closes and removes all pooled connections.
func (p *connPool) closeAll() {
	p.mu.Lock()
	conns := p.conns
	p.conns = make(map[string]*nwep.Client)
	p.mu.Unlock()

	for _, c := range conns {
		c.Close()
	}
}
