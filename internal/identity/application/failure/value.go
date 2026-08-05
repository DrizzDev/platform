package failure

import (
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/failure/action"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/category"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

type Value struct {
	detail      string
	correlation string
	code        code.Code
	retry       time.Duration
}

func (value Value) Code() code.Code {
	return value.code
}

func (value Value) Category() category.Category {
	return value.code.Category()
}

func (value Value) Action() action.Action {
	return value.code.Action()
}

func (value Value) Retryable() bool {
	return value.code.Retryable()
}

func (value Value) Detail() string {
	return value.detail
}

func (value Value) Correlation() string {
	return value.correlation
}

func (value Value) Retry() time.Duration {
	return value.retry
}
