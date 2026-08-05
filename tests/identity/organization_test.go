//go:build cloud

package identity_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/cloud"
)

// TestOrganization exercises the real Drizz Cloud organization endpoint. It runs
// only under the `cloud` build tag and only when the staging base URL and a live
// access token are supplied, so it stays out of the general gate until the
// deployment facts (base URL and desktop-token JWKS domain) are in place. Run:
// `DRIZZ_CLOUD=<base> DRIZZ_ACCESS_TOKEN=<token> go test -tags cloud ./tests/identity`.
func TestOrganization(test *testing.T) {
	base := os.Getenv("DRIZZ_CLOUD")
	token := os.Getenv("DRIZZ_ACCESS_TOKEN")
	if base == "" || token == "" {
		test.Skip("set DRIZZ_CLOUD and DRIZZ_ACCESS_TOKEN to run the staging organization check")
	}
	client, failure := cloud.New(cloud.Options{Agent: &http.Client{Timeout: 15 * time.Second}, Base: base})
	if failure != nil {
		test.Fatal(failure)
	}
	tenant, failure := client.Authorize(context.Background(),
		login.Grant{Token: []byte(token), Expiry: time.Now().Add(time.Hour)})
	if failure != nil {
		test.Fatalf("staging denied the sign-in: %v", failure)
	}
	test.Logf("resolved organization: %q", tenant.Name)
}
