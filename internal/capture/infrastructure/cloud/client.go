package cloud

import (
	"context"
	"net/http"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Client uploads captured artifacts and entries to Drizz Cloud. Each request is idempotent by identity — an artifact
// by its digest, an entry by its span — so the courier's at-least-once resend converges on a single stored effect.
// Retry, backoff, jitter, and cancellation are the courier's concern; timeouts and observability are the transport
// client's; this adapter only builds the request and maps the outcome back to the courier's retry contract.
type Client struct {
	base   string
	source *source
	agent  *http.Client
}

// Blob uploads one artifact's bytes, keyed by its digest. A repeat upload of the same digest is a no-op at the cloud.
func (client Client) Blob(scope context.Context, cargo courier.Cargo) error {
	request, failure := client.stream(scope, cargo)

	if failure != nil {
		return failure
	}
	return client.send(scope, request)
}

// Record uploads one journal entry, keyed by its span. A repeat upload of the same span is a no-op at the cloud.
func (client Client) Record(scope context.Context, entry journal.Entry) error {
	request, failure := client.encode(scope, entry)

	if failure != nil {
		return failure
	}
	return client.send(scope, request)
}

// send authorizes and performs one request, then maps the response to the courier's retry contract. A cancelled or
// timed-out scope surfaces as itself so the courier stops; any other transport fault is transient and gets retried.
func (client Client) send(scope context.Context, request *http.Request) error {
	credential, failure := client.source.acquire(scope)

	if failure != nil {
		return failure
	}
	request.Header.Set("Authorization", "Bearer "+string(credential.Token))
	response, failure := client.agent.Do(request)

	if failure != nil {
		if scope.Err() != nil {
			return scope.Err()
		}
		return unreachable{}
	}
	defer func() { _ = response.Body.Close() }()
	return client.interpret(response)
}
