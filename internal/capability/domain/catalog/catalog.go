package catalog

// The names of the capabilities Drizz offers. They are defined once here so every surface and the code that performs
// them refer to the same name and can never drift apart.
const (
	Devices    = "devices"
	Screenshot = "screenshot"
	Tap        = "tap"
)

// Catalog is the single, ordered list of the capabilities Drizz offers. The command line and the agent connection both
// read from it, so the two always present the same capabilities described the same way.
type Catalog struct {
	entries []Entry
}

func New() Catalog {
	return Catalog{entries: []Entry{
		{
			name:    Devices,
			title:   "ListDevices",
			summary: "List the devices currently connected to this computer.",
		},
		{
			name:    Screenshot,
			title:   "TakeScreenshot",
			summary: "Capture the current screen of a connected device.",
			parameters: []Parameter{
				{name: "serial", summary: "The serial of the device to capture."},
			},
		},
		{
			name:    Tap,
			title:   "Tap",
			summary: "Tap the screen of a connected device at a point.",
			parameters: []Parameter{
				{name: "serial", summary: "The serial of the device to tap."},
				{name: "x", summary: "The horizontal coordinate to tap, in pixels."},
				{name: "y", summary: "The vertical coordinate to tap, in pixels."},
			},
		},
	}}
}

func (catalog Catalog) List() []Entry {
	return append([]Entry(nil), catalog.entries...)
}

func (catalog Catalog) Lookup(name string) (Entry, bool) {
	for _, entry := range catalog.entries {
		if entry.name == name {
			return entry, true
		}
	}
	return Entry{}, false
}
