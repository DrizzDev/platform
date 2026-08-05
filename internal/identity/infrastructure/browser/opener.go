package browser

import "context"

// Opener opens a URL in the user's system browser. It is a port so the loopback
// flow can be driven without a real browser under test.
type Opener interface {
	Open(context.Context, string) error
}
