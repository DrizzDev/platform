//go:build keyring

package identity_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/vault"
)

func TestKeyring(test *testing.T) {
	safe, failure := vault.New(vault.Options{Store: vault.Keyring{}})
	if failure != nil {
		test.Fatal(failure)
	}

	account, failure := subject.New("google-oauth2|identity-keyring-probe")
	if failure != nil {
		test.Fatal(failure)
	}
	handle, failure := session.New("identity-keyring-probe")
	if failure != nil {
		test.Fatal(failure)
	}
	record, failure := credential.New(credential.Input{
		Issuer:   "https://issuer.example/",
		Client:   "native",
		Handle:   "handle_1234567890",
		Subject:  account,
		Session:  handle,
		Method:   method.Browser,
		Refresh:  []byte("probe-refresh-bytes"),
		Issued:   time.Unix(1000, 0),
		Expiry:   time.Unix(2000, 0),
		Revision: 1,
		Schema:   1,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	test.Cleanup(func() {
		_ = safe.Delete(context.Background(), record.Key())
	})

	if failure := safe.Write(context.Background(), record); failure != nil {
		test.Fatal(failure)
	}
	restored, failure := safe.Read(context.Background(), record.Key())
	if failure != nil {
		test.Fatal(failure)
	}
	if !bytes.Equal(restored.Refresh(), record.Refresh()) || restored.Session().String() != handle.String() {
		test.Fatalf("credential did not round-trip through the real store: %+v", restored)
	}

	if failure := safe.Delete(context.Background(), record.Key()); failure != nil {
		test.Fatal(failure)
	}
	if _, absent := safe.Read(context.Background(), record.Key()); absent == nil {
		test.Fatal("a deleted credential is still present")
	}
}
