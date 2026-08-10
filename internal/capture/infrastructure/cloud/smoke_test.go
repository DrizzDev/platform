//go:build cloud

package cloud_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/cloud"
)

// live yields a static bearer for the smoke run; the real composition root supplies an identity-backed provider.
type live struct {
	token string
}

func (fake live) Access(context.Context) (cloud.Credential, error) {
	return cloud.Credential{Token: []byte(fake.token), Expiry: time.Unix(1<<62, 0)}, nil
}

// TestLiveUpload exercises the adapter against a real endpoint. It runs only under `-tags cloud` and only when the
// endpoint and a token are supplied, so it stays out of the offline suite until the cloud route is live.
func TestLiveUpload(test *testing.T) {
	base := os.Getenv("DRIZZ_CLOUD")
	token := os.Getenv("DRIZZ_CLOUD_TOKEN")
	if base == "" || token == "" {
		test.Skip("set DRIZZ_CLOUD and DRIZZ_CLOUD_TOKEN to run the live cloud smoke")
	}

	kit := fixture{test: test}
	client := kit.client(cloud.Options{Provider: live{token: token}, Base: base})

	blob := []byte("smoke-artifact")
	cargo := courier.Cargo{Digest: kit.digest(blob), Source: strings.NewReader(string(blob))}
	if failure := client.Blob(context.Background(), cargo); failure != nil {
		test.Fatalf("blob: %v", failure)
	}
	if failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-smoke", blob: blob})); failure != nil {
		test.Fatalf("record: %v", failure)
	}
}
