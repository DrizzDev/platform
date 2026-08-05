package logout

import "time"

// Clock supplies the current time for the cleanup enqueue timestamps.
type Clock interface {
	Now() time.Time
}
