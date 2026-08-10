package artifact

// Temp exposes the staging directory so a test can plant an orphaned file and prove Sweep removes it.
func (store Store) Temp() string {
	return store.temp
}
