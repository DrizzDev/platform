package binary

// Using builds a Resolver backed by a caller-supplied executable-path function, so a test can drive Locate without
// depending on the real process path.
func Using(executable func() (string, error)) Resolver {
	return Resolver{executable: executable}
}
