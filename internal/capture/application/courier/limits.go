package courier

import "time"

const (
	attempts = 5
	ceiling  = 5 * time.Second
	lull     = 200 * time.Millisecond
)
