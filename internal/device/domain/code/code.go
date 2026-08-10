package code

// Code is the stable, agent-facing device failure vocabulary. Each code owns its safe, actionable detail; no raw
// cause, vendor string, path, endpoint, or serial is ever carried here.
type Code string

const (
	Timeout      Code = "DEVICE_TIMEOUT"
	Missing      Code = "DEVICE_NOT_FOUND"
	Failed       Code = "DEVICE_COMMAND_FAILED"
	Protected    Code = "SCREEN_PROTECTED"
	Cancelled    Code = "DEVICE_CANCELLED"
	Unavailable  Code = "DEVICE_UNAVAILABLE"
	Incompatible Code = "DEVICE_INCOMPATIBLE"
	Unauthorized Code = "DEVICE_UNAUTHORIZED"
)

func (code Code) Valid() bool {
	switch code {
	case Missing, Timeout, Failed, Protected, Cancelled, Unavailable, Incompatible, Unauthorized:
		return true
	default:
		return false
	}
}

func (code Code) String() string {
	return string(code)
}

func (code Code) Retryable() bool {
	switch code {
	case Timeout, Unavailable:
		return true
	case Missing, Failed, Protected, Cancelled, Incompatible, Unauthorized:
		return false
	}
	return false
}

// Detail is the allowlisted explanation an agent may act on.
func (code Code) Detail() string {
	switch code {
	case Missing:
		return "No device matched the identifier; ensure it is connected and authorized."
	case Unauthorized:
		return "The device has not authorized this computer; accept the prompt on the device."
	case Protected:
		return "The current screen is protected and cannot be captured."
	case Timeout:
		return "The device did not respond in time; retry or reconnect it."
	case Incompatible:
		return "The device does not support this operation."
	case Failed:
		return "The operation could not be completed on the device."
	case Unavailable:
		return "The device is temporarily unavailable; retry shortly."
	case Cancelled:
		return "The operation was cancelled."
	default:
		return ""
	}
}
