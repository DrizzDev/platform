package operator

import "github.com/DrizzDev/platform/internal/capability/domain/outcome"

// Refusal is a capability failure carrying the stable, agent-facing code, so a surface renders it with the code's safe
// detail and never a vendor, path, or device string.
type Refusal struct {
	code outcome.Code
}

func (refusal Refusal) Error() string {
	return refusal.code.String()
}

func (refusal Refusal) Code() outcome.Code {
	return refusal.code
}
