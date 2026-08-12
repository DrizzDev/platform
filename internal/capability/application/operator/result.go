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

// Snapshot is a captured screen together with its on-screen element tree.
type Snapshot struct {
	Image     []byte
	Format    string
	Hierarchy string
}

// Tree is a device's on-screen element tree.
type Tree struct {
	Hierarchy string
}

// Extent is a device's screen size in pixels.
type Extent struct {
	Width  int
	Height int
}

// Contact selects where on a device to tap.
type Contact struct {
	Serial string
	X      int
	Y      int
}

// Spot is a coordinate on a device screen, in pixels.
type Spot struct {
	X int
	Y int
}

// Drag is a swipe across a device from one point to another over a duration in milliseconds.
type Drag struct {
	Serial       string
	From         Spot
	To           Spot
	Milliseconds int
}

// Squeeze is a pinch around a centre point, from one radius to another, in pixels.
type Squeeze struct {
	Serial string
	Centre Spot
	From   int
	To     int
}

// Key is a hardware or remote button to press on a device.
type Key struct {
	Serial string
	Button string
}

// Entry is text to type into the focused field on a device.
type Entry struct {
	Serial string
	Text   string
}

// Fix is a location to report on a device, in degrees.
type Fix struct {
	Serial string
	Lat    float64
	Lon    float64
}

// Package is an application package to install: the device and the path to the package file.
type Package struct {
	Serial string
	Path   string
}

// Application names an installed application on a device by its package or bundle identifier.
type Application struct {
	Serial string
	App    string
}

// App is one application read from a device: its identifier, a human name, and a note (its version or process id).
type App struct {
	Id   string
	Name string
	Note string
}

// Listing is a set of applications read from a device.
type Listing struct {
	Apps []App
}

// Report is a single text value read from a device, such as the foreground app or its current link.
type Report struct {
	Text string
}

// Measure is a single numeric value read from a device, such as its free disk space in megabytes.
type Measure struct {
	Value int
}

// Images is the set of emulator image names available to run.
type Images struct {
	Names []string
}

// Image is an emulator image to start, named by its image name.
type Image struct {
	Name string
}

// Ack acknowledges a performed action that returns no data of its own.
type Ack struct{}
