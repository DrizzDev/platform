package device_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/device"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type provider struct {
	token   login.Token
	request error
	await   error
}

func (fake provider) Request(context.Context) (device.Instruction, error) {
	return device.Instruction{Code: "device", User: "WDJB-MJHT", Verification: "https://drizz.dev/activate"}, fake.request
}

func (fake provider) Await(context.Context, device.Instruction) (login.Token, error) {
	return fake.token, fake.await
}

type display struct {
	shown bool
	fail  error
}

func (fake *display) Show(context.Context, device.Instruction) error {
	fake.shown = true
	return fake.fail
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) token() login.Token {
	fixture.test.Helper()
	account, failure := subject.New("google-oauth2|110000000000000000000")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return login.Token{
		Issuer: "https://issuer.example/", Client: "native", Subject: account, Method: method.Device,
		Refresh: []byte("refresh"), Issued: time.Unix(1000, 0), Expiry: time.Unix(2000, 0),
	}
}

func (fixture fixture) establish(options device.Options) (login.Token, error) {
	fixture.test.Helper()
	terminal, failure := device.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return terminal.Establish(context.Background())
}

func TestEstablish(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	screen := &display{}
	token, failure := fixture.establish(device.Options{Provider: provider{token: fixture.token()}, Display: screen})
	if failure != nil {
		test.Fatalf("a valid device grant failed: %v", failure)
	}
	if !screen.shown {
		test.Fatal("the challenge was never shown to the user")
	}
	if token.Subject.String() != "google-oauth2|110000000000000000000" {
		test.Fatalf("subject = %q", token.Subject.String())
	}
}

func TestHalts(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	request := errors.New("request failed")
	shown := errors.New("show failed")
	await := errors.New("await failed")
	cases := map[string]struct {
		options device.Options
		want    error
	}{
		"request": {device.Options{Provider: provider{request: request}, Display: &display{}}, request},
		"display": {device.Options{Provider: provider{token: fixture.token()}, Display: &display{fail: shown}}, shown},
		"await":   {device.Options{Provider: provider{await: await}, Display: &display{}}, await},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			_, failure := fixture.establish(scenario.options)
			if !errors.Is(failure, scenario.want) {
				test.Fatalf("failure = %v, want %v", failure, scenario.want)
			}
		})
	}
}
