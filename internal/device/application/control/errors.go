package control

import "github.com/DrizzDev/platform/internal/device/domain/code"

// The errors below are the control-owned failure contract the bridge adapter returns, so the flow maps an outcome to
// a stable code without inspecting any vendor or sidecar detail. Cancellation is derived from the caller context, not
// carried here.

// Missing reports that no device matched the requested identifier.
type Missing struct{}

func (Missing) Error() string {
	return "no device matched the identifier"
}

func (Missing) reason() code.Code {
	return code.Missing
}

// Unauthorized reports that the device has not authorized this computer.
type Unauthorized struct{}

func (Unauthorized) Error() string {
	return "the device has not authorized this computer"
}

func (Unauthorized) reason() code.Code {
	return code.Unauthorized
}

// Protected reports that the current screen cannot be captured.
type Protected struct{}

func (Protected) Error() string {
	return "the current screen is protected"
}

func (Protected) reason() code.Code {
	return code.Protected
}

// Timeout reports that the device did not respond in time.
type Timeout struct{}

func (Timeout) Error() string {
	return "the device did not respond in time"
}

func (Timeout) reason() code.Code {
	return code.Timeout
}

// Incompatible reports that the device does not support the operation.
type Incompatible struct{}

func (Incompatible) Error() string {
	return "the device does not support the operation"
}

func (Incompatible) reason() code.Code {
	return code.Incompatible
}

// Failed reports that the operation could not be completed on the device.
type Failed struct{}

func (Failed) Error() string {
	return "the operation could not be completed"
}

func (Failed) reason() code.Code {
	return code.Failed
}

// Unavailable reports that the device is temporarily unavailable.
type Unavailable struct{}

func (Unavailable) Error() string {
	return "the device is temporarily unavailable"
}

func (Unavailable) reason() code.Code {
	return code.Unavailable
}
