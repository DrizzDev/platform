package browser_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/browser"
)

type opener struct {
	agent  *http.Client
	target string
	repeat int
}

func (fake opener) Open(scope context.Context, _ string) error {
	if fake.target == "" {
		return nil
	}
	rounds := fake.repeat
	if rounds == 0 {
		rounds = 1
	}
	for range rounds {
		request, failure := http.NewRequestWithContext(scope, http.MethodGet, fake.target, nil)
		if failure != nil {
			return failure
		}
		response, failure := fake.agent.Do(request)
		if failure != nil {
			return failure
		}
		_ = response.Body.Close()
	}
	return nil
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) address() string {
	fixture.test.Helper()
	probe, failure := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	address := probe.Addr().String()
	_ = probe.Close()
	return address
}

func (fixture fixture) build(options browser.Options) browser.Browser {
	fixture.test.Helper()
	made, failure := browser.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestPrompt(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	address := fixture.address()
	fake := opener{agent: &http.Client{}, target: "http://" + address + "/callback?code=code_123&state=state_123"}
	made := fixture.build(browser.Options{Opener: fake, Address: address, Path: "/callback"})

	callback, failure := made.Prompt(context.Background(), login.Redirect{URL: "https://issuer.example/authorize"})
	if failure != nil {
		test.Fatal(failure)
	}
	if callback.Code != "code_123" || callback.State != "state_123" {
		test.Fatalf("callback = %+v", callback)
	}
}

func TestRefused(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	address := fixture.address()
	fake := opener{agent: &http.Client{}, target: "http://" + address + "/callback?state=state_123"}
	made := fixture.build(browser.Options{Opener: fake, Address: address, Path: "/callback"})

	_, failure := made.Prompt(context.Background(), login.Redirect{})
	var rejected login.Rejected
	if !errors.As(failure, &rejected) {
		test.Fatalf("a callback without a code was accepted: %v", failure)
	}
}

func TestCancel(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	address := fixture.address()
	made := fixture.build(browser.Options{Opener: opener{}, Address: address, Path: "/callback"})

	scope, cancel := context.WithCancel(context.Background())
	cancel()
	if _, failure := made.Prompt(scope, login.Redirect{}); !errors.Is(failure, context.Canceled) {
		test.Fatalf("cancellation did not propagate: %v", failure)
	}
}

func TestInterrupted(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	address := fixture.address()
	fake := opener{agent: &http.Client{}, target: "http://" + address + "/health", repeat: 3}
	made := fixture.build(browser.Options{Opener: fake, Address: address, Path: "/callback"})

	_, failure := made.Prompt(context.Background(), login.Redirect{})
	var interrupted browser.Interrupted
	if !errors.As(failure, &interrupted) {
		test.Fatalf("stray requests did not close the listener: %v", failure)
	}
}
