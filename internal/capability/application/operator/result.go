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

// Contact selects where on a device to tap.
type Contact struct {
	Serial string
	X      int
	Y      int
}

// Ack acknowledges a performed action that returns no data of its own.
type Ack struct{}
