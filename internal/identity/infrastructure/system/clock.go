package system

import (
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

var _ login.Clock = Clock{}

// Clock reads the host wall clock. It is injected so login stays deterministic
// under test (CODE-013).
type Clock struct{}

func (Clock) Now() time.Time {
	return time.Now()
}
