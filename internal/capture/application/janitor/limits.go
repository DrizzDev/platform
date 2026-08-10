package janitor

import "time"

// ceiling is the default artifact disk budget; relief is the default low-water mark reclaim frees down to before
// lifting the write gate; grace is the default window that spares a just-written object. Options override each.
const (
	ceiling = 1 << 30
	relief  = 768 << 20
	grace   = time.Hour
)
