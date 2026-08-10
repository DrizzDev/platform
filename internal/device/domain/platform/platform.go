package platform

// Platform is the device family a request is routed to. The values are the internal representation; the bridge
// adapter translates them to and from the sidecar's wire tokens.
type Platform string

const (
	Android   Platform = "ANDROID"
	Handset   Platform = "IOS_DEVICE"
	Simulator Platform = "IOS_SIMULATOR"
)

func (platform Platform) Valid() bool {
	switch platform {
	case Android, Simulator, Handset:
		return true
	default:
		return false
	}
}

func (platform Platform) String() string {
	return string(platform)
}
