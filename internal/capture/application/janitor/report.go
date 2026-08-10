package janitor

// Report is the outcome of one reclaim pass: how many entries it removed and whether the store is now degraded and
// refusing new writes. The caller maps it to telemetry; this package emits none itself.
type Report struct {
	Reclaimed int
	Degraded  bool
}
