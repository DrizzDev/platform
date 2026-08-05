package login

import "time"

// Clock supplies the current time so the flow stays deterministic under test.
type Clock interface {
	Now() time.Time
}
