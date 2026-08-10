package bridge

import "time"

const (
	dialect  = 1
	inflight = 64
	patience = 3        // consecutive request timeouts that recycle a wedged-but-alive sidecar
	frame    = 64 << 20 // a screenshot frame is large; cap well above a base64-encoded image
	grace    = 5 * time.Second
	ceiling  = 10 * time.Second
	lull     = 250 * time.Millisecond
)

const (
	version = "2.0"
	ping    = "health"
	revoke  = "$/cancel"
)
