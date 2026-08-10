package bridge

import (
	"context"
	"encoding/json"
)

type report struct {
	Status   string `json:"status"`
	Version  string `json:"version"`
	Protocol int    `json:"protocol"`
}

// greet handshakes and refuses an incompatible protocol, failing closed rather than
// speaking a version it does not understand.
func (session *session) greet(scope context.Context) error {
	response, failure := session.exchange(scope, Request{Method: ping})
	if failure != nil {
		return failure
	}

	var health report
	if json.Unmarshal(response.Result, &health) != nil {
		return Mismatch{}
	}
	if health.Protocol != dialect {
		return Mismatch{}
	}
	return nil
}
