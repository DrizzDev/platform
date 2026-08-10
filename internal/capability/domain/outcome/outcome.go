package outcome

// Code is the stable, agent-facing failure vocabulary for a capability. It is source-neutral, so a failure from the
// device — and later from any other capability source — maps to the one code a surface renders. Each code owns its
// safe, actionable detail; no vendor, path, or internal cause is ever carried here.
type Code string

const (
	Invalid      Code = "CAPABILITY_INVALID"
	Missing      Code = "DEVICE_NOT_FOUND"
	Unauthorized Code = "DEVICE_UNAUTHORIZED"
	Unavailable  Code = "DEVICE_UNAVAILABLE"
	Timeout      Code = "DEVICE_TIMEOUT"
	Unsupported  Code = "DEVICE_UNSUPPORTED"
	Refused      Code = "DEVICE_REFUSED"
	Failed       Code = "CAPABILITY_FAILED"
	Cancelled    Code = "CAPABILITY_CANCELLED"
)

func (code Code) Valid() bool {
	switch code {
	case Invalid, Missing, Unauthorized, Unavailable, Timeout, Unsupported, Refused, Failed, Cancelled:
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
	case Unavailable, Timeout:
		return true
	case Invalid, Missing, Unauthorized, Unsupported, Refused, Failed, Cancelled:
		return false
	}
	return false
}

// Detail is the allowlisted, agent-actionable explanation for a code.
func (code Code) Detail() string {
	switch code {
	case Invalid:
		return "The request was missing or invalid; check the inputs and try again."
	case Missing:
		return "The requested device was not found; ensure it is connected and authorized."
	case Unauthorized:
		return "The device has not authorized this computer; accept the prompt on the device."
	case Unavailable:
		return "The device is temporarily unavailable; try again shortly."
	case Timeout:
		return "The device did not respond in time; retry or reconnect it."
	case Unsupported:
		return "The device does not support this capability."
	case Refused:
		return "The device refused the capability; the current screen may be protected."
	case Cancelled:
		return "The capability was cancelled."
	case Failed:
		return "The capability could not be completed."
	}
	return "The capability could not be completed."
}
