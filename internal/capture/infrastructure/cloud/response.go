package cloud

import (
	"net/http"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// wire is the JSON shape of an uploaded entry. Field names match the journal's columns so a record keeps one
// vocabulary across storage and transport. The contract is provisional until the cloud endpoint is finalized.
type wire struct {
	Trace    string `json:"trace"`
	Span     string `json:"span"`
	Mark     string `json:"mark"`
	State    string `json:"state"`
	Parent   string `json:"parent"`
	Origin   string `json:"origin"`
	Payload  []byte `json:"payload"`
	Artifact string `json:"artifact"`
	Fidelity string `json:"fidelity"`
	Sequence int64  `json:"sequence"`
	Category string `json:"category"`
}

func (Client) compose(entry journal.Entry) wire {
	link := entry.Correlation()

	return wire{
		Sequence: link.Sequence(),
		Payload:  entry.Payload(),
		Mark:     link.Mark().String(),
		Span:     link.Span().String(),
		Trace:    link.Trace().String(),
		Parent:   link.Parent().String(),
		State:    entry.State().String(),
		Origin:   entry.Origin().String(),
		Fidelity: entry.Fidelity().String(),
		Category: entry.Category().String(),
		Artifact: entry.Artifact().String(),
	}
}

// interpret maps the cloud response to the courier's retry contract: success and an idempotent conflict are done;
// a client rejection is terminal and must not be retried; a server fault is transient and will be retried.
func (Client) interpret(response *http.Response) error {
	switch {
	case response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices:
		return nil
	case response.StatusCode == http.StatusConflict:
		return nil
	case response.StatusCode >= http.StatusBadRequest && response.StatusCode < http.StatusInternalServerError:
		return courier.Rejected{}
	default:
		return transient{status: response.StatusCode}
	}
}
