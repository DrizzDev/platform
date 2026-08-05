package failure

import (
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

type Input struct {
	Detail      string
	Correlation string
	Code        code.Code
	Retry       time.Duration
}
