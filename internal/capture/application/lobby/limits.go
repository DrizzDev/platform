package lobby

import "time"

// window is the default lifetime of a pending observation and the default horizon for an inferred time match; capacity
// is the default cap on how many observations a single match considers. Options override each.
const (
	capacity = 256
	window   = 5 * time.Minute
)
