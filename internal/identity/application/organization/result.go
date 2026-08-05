package organization

import (
	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/organization"
)

// Result is the resolved organization context, or a code-only failure. It never
// carries a token, role, permission set, or raw cloud response.
type Result struct {
	failure      *fault.Value
	organization organization.Organization
}

func (result Result) Organization() organization.Organization {
	return result.organization
}

func (result Result) Failed() bool {
	return result.failure != nil
}

func (result Result) Failure() (fault.Value, bool) {
	if result.failure == nil {
		return fault.Value{}, false
	}
	return *result.failure, true
}
