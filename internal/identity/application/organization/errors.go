package organization

import "github.com/DrizzDev/platform/internal/identity/domain/failure/code"

// The port errors below are the organization-owned failure contract. The cloud
// adapter returns them so the flow maps an outcome to a stable code without
// inspecting the cloud response.

// Forbidden reports that the subject has no active membership in the resolved
// organization.
type Forbidden struct{}

func (Forbidden) Error() string {
	return "the subject is not permitted in the organization"
}

func (Forbidden) reason() code.Code {
	return code.Forbidden
}

// Required reports that no organization can be resolved for the subject and a
// new sign-in is needed.
type Required struct{}

func (Required) Error() string {
	return "a new sign-in is required to resolve the organization"
}

func (Required) reason() code.Code {
	return code.Required
}

// Unavailable reports that Drizz Cloud could not be reached or returned an
// unusable response.
type Unavailable struct{}

func (Unavailable) Error() string {
	return "organization resolution is temporarily unavailable"
}

func (Unavailable) reason() code.Code {
	return code.Unavailable
}
