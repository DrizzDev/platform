package operator

import (
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/domain/code"
)

// refuse maps a device failure code into the capability's own agent-facing code, so a surface depends only on the
// capability vocabulary and a future non-device source can map into the same set.
func (Operator) refuse(cause code.Code) Refusal {
	switch cause {
	case code.Missing:
		return Refusal{code: outcome.Missing}
	case code.Unauthorized:
		return Refusal{code: outcome.Unauthorized}
	case code.Unavailable:
		return Refusal{code: outcome.Unavailable}
	case code.Timeout:
		return Refusal{code: outcome.Timeout}
	case code.Incompatible:
		return Refusal{code: outcome.Unsupported}
	case code.Protected:
		return Refusal{code: outcome.Refused}
	case code.Cancelled:
		return Refusal{code: outcome.Cancelled}
	case code.Failed:
		return Refusal{code: outcome.Failed}
	default:
		return Refusal{code: outcome.Failed}
	}
}
