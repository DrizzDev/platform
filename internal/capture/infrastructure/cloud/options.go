package cloud

import "net/http"

// Options configures the cloud uploader. Agent is the resilient, observable client from platform/transport — it owns
// the timeouts, connection bounds, response-body cap, and OpenTelemetry spans, so this adapter adds none of its own.
type Options struct {
	Clock    Clock
	Base     string
	Provider Provider
	Agent    *http.Client
}
