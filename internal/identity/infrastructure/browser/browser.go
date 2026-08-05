package browser

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

const (
	timeout = 5 * time.Second
	ceiling = 8 << 10
)

var _ login.Browser = Browser{}

// Browser implements the login Browser port. It binds a loopback listener,
// opens the system browser, and returns the single captured callback.
type Browser struct {
	opener  Opener
	address string
	path    string
}

func (browser Browser) Prompt(scope context.Context, redirect login.Redirect) (login.Callback, error) {
	listener, failure := (&net.ListenConfig{}).Listen(scope, "tcp", browser.address)
	if failure != nil {
		return login.Callback{}, failure
	}
	relay := &capture{callbacks: make(chan login.Callback, 1), faults: make(chan error, 1)}
	server := &http.Server{Handler: browser.route(relay), ReadHeaderTimeout: timeout, MaxHeaderBytes: ceiling}
	go func() { _ = server.Serve(listener) }()
	defer browser.shutdown(scope, server)

	if failure := browser.opener.Open(scope, redirect.URL); failure != nil {
		return login.Callback{}, failure
	}
	select {
	case <-scope.Done():
		return login.Callback{}, scope.Err()
	case failure := <-relay.faults:
		return login.Callback{}, failure
	case callback := <-relay.callbacks:
		return callback, nil
	}
}

func (browser Browser) shutdown(scope context.Context, server *http.Server) {
	stopping, cancel := context.WithTimeout(context.WithoutCancel(scope), timeout)
	defer cancel()
	_ = server.Shutdown(stopping)
}
