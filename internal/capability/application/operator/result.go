package operator

// Target selects the device a capability acts on.
type Target struct {
	Serial string
}

// Shot is a captured screen: the image bytes and the format they are encoded in.
type Shot struct {
	Image  []byte
	Format string
}

// Roster is the set of connected device serials, in discovery order.
type Roster struct {
	Serials []string
}
