package cloud_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/cloud"
)

type provider struct {
	credential cloud.Credential
	fail       error
	calls      int
}

func (fake *provider) Access(context.Context) (cloud.Credential, error) {
	fake.calls++
	return fake.credential, fake.fail
}

type clock struct{}

func (clock) Now() time.Time { return time.Unix(1000, 0) }

type spec struct {
	hop  string
	blob []byte
}

// recorder is a fake cloud endpoint that answers one fixed status and keeps the single request the adapter sent.
type recorder struct {
	server *httptest.Server
	seen   http.Request
	body   []byte
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) credential() cloud.Credential {
	return cloud.Credential{Token: []byte("access-token"), Expiry: time.Unix(9000, 0)}
}

func (fixture fixture) serve(status int) *recorder {
	record := &recorder{}
	record.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		record.body = body
		record.seen = *request
		writer.WriteHeader(status)
	}))
	return record
}

func (fixture fixture) client(options cloud.Options) cloud.Client {
	fixture.test.Helper()
	options.Agent = &http.Client{}
	options.Clock = clock{}
	made, failure := cloud.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) digest(blob []byte) digest.Digest {
	fixture.test.Helper()
	sum := sha256.Sum256(blob)
	made, failure := digest.New(hex.EncodeToString(sum[:]))
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) entry(spec spec) journal.Entry {
	fixture.test.Helper()
	thread, failure := trace.New("01HEXECUTION")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	point, failure := span.New(spec.hop)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	link, failure := correlation.New(correlation.Input{Trace: thread, Span: point, Mark: mark.Exact})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	reference := digest.Digest{}
	if len(spec.blob) > 0 {
		reference = fixture.digest(spec.blob)
	}
	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Origin:      origin.Capability,
		Fidelity:    fidelity.Exact,
		Category:    category.Tool,
		Artifact:    reference,
		State:       journal.Pending,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return entry
}

func TestRecord(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusOK)
	defer cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	if failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"})); failure != nil {
		test.Fatalf("record: %v", failure)
	}

	if cloudside.seen.Method != http.MethodPut || cloudside.seen.URL.Path != "/captures/entries/hop-1" {
		test.Fatalf("request = %s %s", cloudside.seen.Method, cloudside.seen.URL.Path)
	}
	if cloudside.seen.Header.Get("Authorization") != "Bearer access-token" {
		test.Fatalf("authorization = %q", cloudside.seen.Header.Get("Authorization"))
	}
	var payload struct {
		Span     string `json:"span"`
		Category string `json:"category"`
	}
	if failure := json.Unmarshal(cloudside.body, &payload); failure != nil {
		test.Fatalf("body: %v", failure)
	}
	if payload.Span != "hop-1" || payload.Category == "" {
		test.Fatalf("payload = %+v", payload)
	}
}

func TestBlob(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	blob := []byte("artifact-bytes")
	cloudside := kit.serve(http.StatusOK)
	defer cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	cargo := courier.Cargo{Digest: kit.digest(blob), Source: strings.NewReader(string(blob))}
	if failure := client.Blob(context.Background(), cargo); failure != nil {
		test.Fatalf("blob: %v", failure)
	}

	if cloudside.seen.Method != http.MethodPut || !strings.HasPrefix(cloudside.seen.URL.Path, "/captures/artifacts/") {
		test.Fatalf("request = %s %s", cloudside.seen.Method, cloudside.seen.URL.Path)
	}
	if string(cloudside.body) != string(blob) {
		test.Fatalf("uploaded body = %q", cloudside.body)
	}
}

func TestConflict(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusConflict)
	defer cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	if failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"})); failure != nil {
		test.Fatalf("conflict must be an idempotent success, got %v", failure)
	}
}

func TestRejected(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusBadRequest)
	defer cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"}))
	var rejected courier.Rejected
	if !errors.As(failure, &rejected) {
		test.Fatalf("a client rejection must be terminal, got %v", failure)
	}
}

func TestTransient(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusServiceUnavailable)
	defer cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"}))
	var rejected courier.Rejected
	if failure == nil || errors.As(failure, &rejected) {
		test.Fatalf("a server fault must be retryable, got %v", failure)
	}
}

func TestUnreachable(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusOK)
	base := cloudside.server.URL
	cloudside.server.Close()

	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: base})
	failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"}))
	var rejected courier.Rejected
	if failure == nil || errors.As(failure, &rejected) {
		test.Fatalf("an unreachable cloud must be retryable, got %v", failure)
	}
}

func TestCancelled(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusOK)
	defer cloudside.server.Close()

	scope, halt := context.WithCancel(context.Background())
	halt()
	client := kit.client(cloud.Options{Provider: &provider{credential: kit.credential()}, Base: cloudside.server.URL})
	if failure := client.Record(scope, kit.entry(spec{hop: "hop-1"})); !errors.Is(failure, context.Canceled) {
		test.Fatalf("a cancelled scope must surface as itself, got %v", failure)
	}
}

func TestCaches(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cloudside := kit.serve(http.StatusOK)
	defer cloudside.server.Close()

	source := &provider{credential: kit.credential()}
	client := kit.client(cloud.Options{Provider: source, Base: cloudside.server.URL})
	for range 3 {
		if failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"})); failure != nil {
			test.Fatalf("record: %v", failure)
		}
	}
	if source.calls != 1 {
		test.Fatalf("credential acquired %d times, want 1 while unexpired", source.calls)
	}
}

func TestProviderFailure(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	broken := &provider{fail: errors.New("no session")}
	client := kit.client(cloud.Options{Provider: broken, Base: "http://cloud.invalid"})
	if failure := client.Record(context.Background(), kit.entry(spec{hop: "hop-1"})); failure == nil {
		test.Fatal("a credential failure must surface")
	}
}
