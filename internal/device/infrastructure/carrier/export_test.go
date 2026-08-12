package carrier

// Probe is arbitrary helper bytes and a digest, so an external test can drive the materialization logic without the
// real embedded helper.
type Probe struct {
	Bytes  []byte
	Digest string
}

// Wrap builds a Carrier rooted at a test directory.
func Wrap(root string) Carrier {
	return Carrier{root: root}
}

// Deliver runs the materialization logic against caller-supplied bytes and digest.
func (carrier Carrier) Deliver(probe Probe) (string, error) {
	return carrier.deliver(asset{bytes: probe.Bytes, digest: probe.Digest})
}
