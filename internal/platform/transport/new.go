package transport

import (
	"context"
	"net"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// New builds a resilient, observable HTTP client: bounded connect and
// per-attempt timeouts, a capped response body, OpenTelemetry spans and metrics,
// and jittered-backoff retries restricted to idempotent requests so a rotating
// secret is never resubmitted after a lost response.
func New(options Options) (*http.Client, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	base := &http.Transport{DialContext: (&net.Dialer{Timeout: options.Dial}).DialContext}
	traced := otelhttp.NewTransport(
		limiter{inner: base, ceiling: options.Ceiling},
		otelhttp.WithTracerProvider(options.Tracing),
		otelhttp.WithMeterProvider(options.Metering),
		otelhttp.WithPropagators(options.propagator()),
	)
	client := retryablehttp.NewClient()
	client.HTTPClient = &http.Client{Transport: traced, Timeout: options.Timeout}
	client.RetryMax = options.Retries
	client.RetryWaitMin = options.Minimum
	client.RetryWaitMax = options.Maximum
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.Logger = nil
	client.CheckRetry = func(_ context.Context, response *http.Response, failure error) (bool, error) {
		if failure != nil || response == nil || response.Request == nil {
			return false, failure
		}
		if response.Request.Method != http.MethodGet && response.Request.Method != http.MethodHead {
			return false, nil
		}
		switch response.StatusCode {
		case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true, nil
		default:
			return false, nil
		}
	}
	return client.StandardClient(), nil
}
