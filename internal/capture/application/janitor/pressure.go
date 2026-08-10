package janitor

// Pressure reports that the artifact store is over its disk budget after reclaim: the bytes in use and the ceiling
// they crossed. It carries only bounded scalars, so it is safe to surface to the user and to turn into a metric.
type Pressure struct {
	Used    int64
	Ceiling int64
}
