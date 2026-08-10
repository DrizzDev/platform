package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

const (
	entries   = "/captures/entries/"
	artifacts = "/captures/artifacts/"
)

// stream builds an idempotent artifact upload: a PUT to the digest-addressed route carrying the raw bytes.
func (client Client) stream(scope context.Context, cargo courier.Cargo) (*http.Request, error) {
	route := client.base + artifacts + cargo.Digest.String()

	request, failure := http.NewRequestWithContext(scope, http.MethodPut, route, cargo.Source)
	if failure != nil {
		return nil, failure
	}
	request.Header.Set("Content-Type", "application/octet-stream")

	return request, nil
}

// encode builds an idempotent entry upload: a PUT to the span-addressed route carrying the entry as JSON.
func (client Client) encode(scope context.Context, entry journal.Entry) (*http.Request, error) {
	body, failure := json.Marshal(client.compose(entry))
	if failure != nil {
		return nil, failure
	}

	route := client.base + entries + entry.Correlation().Span().String()
	request, failure := http.NewRequestWithContext(scope, http.MethodPut, route, bytes.NewReader(body))
	if failure != nil {
		return nil, failure
	}

	request.Header.Set("Content-Type", "application/json")
	return request, nil
}
