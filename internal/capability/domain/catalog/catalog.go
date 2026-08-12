package catalog

// The names of the capabilities Drizz offers. They are defined once here so every surface and the code that performs
// them refer to the same name and can never drift apart.
const (
	Devices    = "devices"
	Screenshot = "screenshot"
	Tap        = "tap"
	Swipe      = "swipe"
	Pinch      = "pinch"
	Press      = "press"
	Type       = "type"
	Clear      = "clear"
	Back       = "back"
	Home       = "home"
	Locate     = "locate"
	Snapshot   = "snapshot"
	Hierarchy  = "hierarchy"
	Dimensions = "dimensions"
	Install    = "install"
	Launch     = "launch"
	Terminate  = "terminate"
	Wipe       = "wipe"
	Installed  = "installed"
	Running    = "running"
	Foreground = "foreground"
	Url        = "url"
	Disk       = "disk"
	Images     = "images"
	Boot       = "boot"
	Pause      = "pause"
	Resume     = "resume"
)

// Catalog is the single, ordered list of the capabilities Drizz offers. The command line and the agent connection both
// read from it, so the two always present the same capabilities described the same way.
type Catalog struct {
	entries []Entry
}

func New() Catalog {
	return Catalog{entries: catalogue}
}

// catalogue is the canonical, ordered set of capabilities. It is read-only: List copies before returning it.
var catalogue = []Entry{
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
	{
		name:    Swipe,
		title:   "Swipe",
		summary: "Swipe across the screen of a connected device from one point to another.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device to swipe."},
			{name: "from-x", summary: "The horizontal coordinate to start from, in pixels."},
			{name: "from-y", summary: "The vertical coordinate to start from, in pixels."},
			{name: "to-x", summary: "The horizontal coordinate to end at, in pixels."},
			{name: "to-y", summary: "The vertical coordinate to end at, in pixels."},
		},
	},
	{
		name:    Pinch,
		title:   "Pinch",
		summary: "Pinch the screen of a connected device to zoom in or out around a point.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device to pinch."},
			{name: "center-x", summary: "The horizontal coordinate of the pinch centre, in pixels."},
			{name: "center-y", summary: "The vertical coordinate of the pinch centre, in pixels."},
			{name: "start-radius", summary: "The starting radius from the centre, in pixels."},
			{name: "end-radius", summary: "The ending radius from the centre, in pixels."},
		},
	},
	{
		name:    Press,
		title:   "PressButton",
		summary: "Press a hardware or remote button on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "button", summary: "The button to press, such as up, down, select, menu, or home."},
		},
	},
	{
		name:    Type,
		title:   "TypeText",
		summary: "Type text into the focused field on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "text", summary: "The text to type."},
		},
	},
	{
		name:    Clear,
		title:   "ClearText",
		summary: "Clear the focused text field on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Back,
		title:   "PressBack",
		summary: "Press the back button on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Home,
		title:   "PressHome",
		summary: "Press the home button on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Locate,
		title:   "SetLocation",
		summary: "Set the reported location of a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "latitude", summary: "The latitude to report, in degrees."},
			{name: "longitude", summary: "The longitude to report, in degrees."},
		},
	},
	{
		name:    Snapshot,
		title:   "TakeSnapshot",
		summary: "Capture the screen of a connected device together with its on-screen element tree.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device to capture."},
		},
	},
	{
		name:    Hierarchy,
		title:   "GetUIHierarchy",
		summary: "Read the on-screen element tree of a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Dimensions,
		title:   "GetScreenSize",
		summary: "Read the screen size of a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Install,
		title:   "InstallApp",
		summary: "Install an application package on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "path", summary: "The path to the application package file."},
		},
	},
	{
		name:    Launch,
		title:   "LaunchApp",
		summary: "Launch an application on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "app", summary: "The application's package or bundle identifier."},
		},
	},
	{
		name:    Terminate,
		title:   "TerminateApp",
		summary: "Stop a running application on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "app", summary: "The application's package or bundle identifier."},
		},
	},
	{
		name:    Wipe,
		title:   "ClearAppData",
		summary: "Clear an application's stored data on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
			{name: "app", summary: "The application's package or bundle identifier."},
		},
	},
	{
		name:    Installed,
		title:   "ListInstalledApps",
		summary: "List the applications installed on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Running,
		title:   "ListRunningApps",
		summary: "List the applications currently running on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Foreground,
		title:   "GetForegroundApp",
		summary: "Read which application is in the foreground on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Url,
		title:   "GetCurrentURL",
		summary: "Read the current link open in the active application on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Disk,
		title:   "GetFreeDiskSpace",
		summary: "Read the free disk space on a connected device.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the device."},
		},
	},
	{
		name:    Images,
		title:   "ListEmulatorImages",
		summary: "List the emulator images available to run on this computer.",
	},
	{
		name:    Boot,
		title:   "BootEmulator",
		summary: "Start an emulator from an available image.",
		parameters: []Parameter{
			{name: "image", summary: "The name of the emulator image to start."},
		},
	},
	{
		name:    Pause,
		title:   "PauseEmulator",
		summary: "Pause a running emulator.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the running emulator."},
		},
	},
	{
		name:    Resume,
		title:   "ResumeEmulator",
		summary: "Resume a paused emulator.",
		parameters: []Parameter{
			{name: "serial", summary: "The serial of the paused emulator."},
		},
	},
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
