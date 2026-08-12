package carrier

import _ "embed"

// bridge and fingerprint are the compiled device helper and its pinned digest, carried inside the platform binary.
// The committed files under asset/ are a development placeholder; the release build replaces them with the compiled
// helper for the target and its regenerated digest before compiling.

//go:embed asset/bridge
var bridge []byte

//go:embed asset/bridge.sha256
var fingerprint string
