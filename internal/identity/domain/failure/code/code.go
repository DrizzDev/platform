package code

import (
	"github.com/DrizzDev/platform/internal/identity/domain/failure/action"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/category"
)

type Code string

const (
	Partial     Code = "LOGOUT_PARTIAL"
	Failed      Code = "IDENTITY_FAILED"
	Conflict    Code = "ACCOUNT_CONFLICT"
	Forbidden   Code = "AUTHORIZATION_FORBIDDEN"
	Rejected    Code = "AUTHENTICATION_REJECTED"
	Required    Code = "AUTHENTICATION_REQUIRED"
	Cancelled   Code = "AUTHENTICATION_CANCELLED"
	Unavailable Code = "AUTHENTICATION_UNAVAILABLE"
	Storage     Code = "SECURE_STORAGE_UNAVAILABLE"
)

func (code Code) Valid() bool {
	switch code {
	case Required, Cancelled, Forbidden, Unavailable, Rejected, Conflict, Storage, Partial, Failed:
		return true
	default:
		return false
	}
}

func (code Code) String() string {
	return string(code)
}

func (code Code) Category() category.Category {
	switch code {
	case Required, Cancelled, Unavailable, Rejected:
		return category.Authentication
	case Forbidden:
		return category.Authorization
	case Conflict:
		return category.Account
	case Storage:
		return category.Storage
	case Partial:
		return category.Logout
	case Failed:
		return category.Internal
	}
	return category.Internal
}

func (code Code) Action() action.Action {
	switch code {
	case Required:
		return action.Signin
	case Conflict:
		return action.Logout
	case Unavailable, Rejected, Failed:
		return action.Retry
	case Cancelled, Forbidden, Storage, Partial:
		return action.None
	}
	return action.None
}

func (code Code) Retryable() bool {
	switch code {
	case Unavailable, Rejected, Failed:
		return true
	case Required, Cancelled, Forbidden, Conflict, Storage, Partial:
		return false
	}
	return false
}
