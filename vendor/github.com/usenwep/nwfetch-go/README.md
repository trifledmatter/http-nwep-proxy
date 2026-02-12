# nwfetch

A WEB/1 fetch client for Go, built on [nwep-go](https://github.com/usenwep/nwep-go).

nwfetch handles connection pooling, identity management, and request construction so you can talk to WEB/1 servers with a familiar HTTP-style API. Every connection is mutually authenticated with Ed25519, and nwfetch manages ephemeral or persistent identities transparently.

nwep-go types (`nwep.Keypair`, `nwep.Header`, `nwep.Notification`, `nwep.Settings`) are used directly. Nothing is re-wrapped.

## Installation

Requires __Go 1.25__ or higher.

```sh
go get github.com/usenwep/nwfetch-go
```

nwfetch depends on [nwep-go](https://github.com/usenwep/nwep-go), which needs platform-specific C libraries. After adding nwfetch to your project, vendor your dependencies and run the nwep-go setup script:

```sh
go mod vendor
cd vendor/github.com/usenwep/nwep-go && bash setup.sh
go build -mod=vendor ./...
```

The setup script downloads pre-built nwep binaries for your platform. It only needs to run once (or again after updating dependencies).

### Makefile for your project

To avoid repeating these steps, drop this into your project's Makefile:

```makefile
NWEP_VENDOR := vendor/github.com/usenwep/nwep-go
STAMP := $(NWEP_VENDOR)/.nwep-setup

build: $(STAMP)
	go build -mod=vendor ./...

test: $(STAMP)
	go test -mod=vendor ./...

$(STAMP): go.mod go.sum
	go mod vendor
	cd $(NWEP_VENDOR) && bash setup.sh
	@touch $@

clean:
	rm -rf vendor
```

Then `make build` handles everything. The setup is cached and only re-runs when `go.mod` or `go.sum` change.

## Quick start

```go
package main

import (
    "fmt"
    "log"

    nwfetch "github.com/usenwep/nwfetch-go"
)

func main() {
    if err := nwfetch.Init(); err != nil {
        log.Fatal(err)
    }
    defer nwfetch.Close()

    resp, err := nwfetch.Get("web://addr/hello")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.String())
}
```

## Request builder

The fluent builder lets you construct requests with method, headers, body, and timeout before executing.

```go
resp, err := nwfetch.New("web://addr/path").
    Method(nwfetch.MethodWrite).
    Header("content-type", "application/json").
    Body([]byte(`{"name": "alice"}`)).
    Do()
```

## Client with options

For long-lived programs or programs that need a stable identity, create a Client directly.

```go
client, err := nwfetch.NewClient(
    nwfetch.WithKeypair(kp),
    nwfetch.WithTimeout(10 * time.Second),
    nwfetch.WithOnNotify(func(n *nwep.Notification) {
        fmt.Printf("notification: %s %s\n", n.Event, n.Path)
    }),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

resp, err := client.Get("web://addr/hello")
```

If no keypair is provided, an ephemeral Ed25519 identity is generated automatically.

## Error handling

Transport errors (connection failures, timeouts) are returned as `*nwfetch.Error` with an Op field indicating the failure stage.

Protocol errors (server returned an error status) can be checked on the response:

```go
resp, err := client.Get("web://addr/resource")
if err != nil {
    log.Fatal(err) // transport error
}
if err := resp.StatusError(); err != nil {
    log.Fatal(err) // protocol error (not_found, forbidden, etc.)
}
```

Sentinel checks work with both `resp.StatusError()` and any wrapped error:

```go
if nwfetch.IsNotFound(err) { ... }
if nwfetch.IsRateLimited(err) {
    if d, ok := resp.RetryAfter(); ok {
        time.Sleep(d)
    }
}
```

## Related

- [nwep-go](https://github.com/usenwep/nwep-go) — Go bindings for libnwep
- [velocity](https://github.com/usenwep/velocity) — WEB/1 server framework for Go

## License

[MIT](LICENSE)
