package auth0_test

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/auth0"
)

const (
	nonce    = "nonce_1234567890"
	state    = "state_1234567890"
	verifier = "verifier_1234567890_1234567890_1234567890_12"
	client   = "platform_cli_client"
	account  = "google-oauth2|110000000000000000000"
)

type fixture struct {
	test     *testing.T
	key      *rsa.PrivateKey
	issuer   string
	audience string
	barren   bool
	stale    bool
}

func (fixture *fixture) start() *httptest.Server {
	fixture.test.Helper()
	key, failure := rsa.GenerateKey(rand.Reader, 2048)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.key = key
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	fixture.issuer = server.URL
	fixture.test.Cleanup(server.Close)

	mux.HandleFunc("/.well-known/openid-configuration", func(writer http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(writer).Encode(map[string]string{
			"issuer":                 fixture.issuer,
			"authorization_endpoint": fixture.issuer + "/authorize",
			"token_endpoint":         fixture.issuer + "/oauth/token",
			"jwks_uri":               fixture.issuer + "/jwks",
		})
	})
	mux.HandleFunc("/jwks", func(writer http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(writer).Encode(fixture.keys())
	})
	mux.HandleFunc("/oauth/token", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		body := map[string]any{
			"access_token": "access-token",
			"token_type":   "Bearer",
			"expires_in":   900,
			"id_token":     fixture.sign(),
		}
		if !fixture.barren {
			body["refresh_token"] = "refresh-token-bytes"
		}
		_ = json.NewEncoder(writer).Encode(body)
	})
	return server
}

func (fixture *fixture) sign() string {
	fixture.test.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT","kid":"test"}`))
	aud := client
	if fixture.audience != "" {
		aud = fixture.audience
	}
	expiry := time.Now().Add(time.Hour)
	if fixture.stale {
		expiry = time.Now().Add(-time.Hour)
	}
	body, failure := json.Marshal(map[string]any{
		"iss":   fixture.issuer,
		"aud":   aud,
		"sub":   account,
		"nonce": nonce,
		"iat":   time.Now().Add(-time.Minute).Unix(),
		"exp":   expiry.Unix(),
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	signing := header + "." + base64.RawURLEncoding.EncodeToString(body)
	digest := sha256.Sum256([]byte(signing))
	signature, failure := rsa.SignPKCS1v15(rand.Reader, fixture.key, crypto.SHA256, digest[:])
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return signing + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func (fixture *fixture) keys() map[string]any {
	public := fixture.key.PublicKey
	return map[string]any{"keys": []map[string]any{{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": "test",
		"n":   base64.RawURLEncoding.EncodeToString(public.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(public.E)).Bytes()),
	}}}
}

func (fixture *fixture) authorizer(agent *http.Client) auth0.Authorizer {
	fixture.test.Helper()
	made, failure := auth0.New(context.Background(), auth0.Options{
		Agent:    agent,
		Issuer:   fixture.issuer,
		Client:   client,
		Audience: "https://platform.drizz.dev",
		Redirect: "http://127.0.0.1:8490/callback",
		Method:   method.Browser,
		Scopes:   []string{"openid", "offline_access"},
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestBegin(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test}
	server := fixture.start()
	redirect, failure := fixture.authorizer(server.Client()).Begin(
		context.Background(), login.Secret{State: state, Nonce: nonce, Verifier: verifier})
	if failure != nil {
		test.Fatal(failure)
	}
	parsed, failure := url.Parse(redirect.URL)
	if failure != nil {
		test.Fatal(failure)
	}
	query := parsed.Query()
	if query.Get("code_challenge_method") != "S256" || query.Get("code_challenge") == "" {
		test.Fatalf("missing PKCE challenge: %s", redirect.URL)
	}
	if query.Get("state") != state || query.Get("nonce") != nonce {
		test.Fatalf("state or nonce not bound: %s", redirect.URL)
	}
	if query.Get("audience") != "https://platform.drizz.dev" {
		test.Fatalf("audience not requested: %s", redirect.URL)
	}
}

func TestFinish(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test}
	server := fixture.start()
	token, failure := fixture.authorizer(server.Client()).Finish(context.Background(), login.Exchange{
		Secret:   login.Secret{State: state, Nonce: nonce, Verifier: verifier},
		Callback: login.Callback{Code: "code_1234567890", State: state},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if token.Subject.String() != account || token.Method != method.Browser {
		test.Fatalf("token = %+v", token)
	}
	if string(token.Refresh) != "refresh-token-bytes" || !token.Expiry.After(token.Issued) {
		test.Fatalf("token lifecycle = %+v", token)
	}
}

func TestReject(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test}
	server := fixture.start()
	_, failure := fixture.authorizer(server.Client()).Finish(context.Background(), login.Exchange{
		Secret:   login.Secret{State: state, Nonce: "wrong_nonce_value", Verifier: verifier},
		Callback: login.Callback{Code: "code_1234567890", State: state},
	})
	var rejected login.Rejected
	if !errors.As(failure, &rejected) {
		test.Fatalf("a mismatched nonce was accepted: %v", failure)
	}
}

func (fixture *fixture) exchange() login.Exchange {
	return login.Exchange{
		Secret:   login.Secret{State: state, Nonce: nonce, Verifier: verifier},
		Callback: login.Callback{Code: "code_1234567890", State: state},
	}
}

func (fixture *fixture) refuse(exchange login.Exchange) error {
	fixture.test.Helper()
	server := fixture.start()
	_, failure := fixture.authorizer(server.Client()).Finish(context.Background(), exchange)
	return failure
}

func TestState(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test}
	failure := fixture.refuse(login.Exchange{
		Secret:   login.Secret{State: state, Nonce: nonce, Verifier: verifier},
		Callback: login.Callback{Code: "code_1234567890", State: "mismatched_state"},
	})
	var rejected login.Rejected
	if !errors.As(failure, &rejected) {
		test.Fatalf("a mismatched callback state was accepted: %v", failure)
	}
}

func TestRefresh(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test, barren: true}
	var rejected login.Rejected
	if failure := fixture.refuse(fixture.exchange()); !errors.As(failure, &rejected) {
		test.Fatalf("a response without a refresh token was accepted: %v", failure)
	}
}

func TestStale(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test, stale: true}
	var rejected login.Rejected
	if failure := fixture.refuse(fixture.exchange()); !errors.As(failure, &rejected) {
		test.Fatalf("an expired identity token was accepted: %v", failure)
	}
}

func TestAudience(test *testing.T) {
	test.Parallel()

	fixture := &fixture{test: test, audience: "https://wrong.example/"}
	var rejected login.Rejected
	if failure := fixture.refuse(fixture.exchange()); !errors.As(failure, &rejected) {
		test.Fatalf("an identity token for a different audience was accepted: %v", failure)
	}
}
