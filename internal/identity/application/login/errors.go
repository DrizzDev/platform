package login

import "github.com/DrizzDev/platform/internal/identity/domain/failure/code"

// The port errors below are the login-owned failure contract. Infrastructure
// returns them so the flow can classify an outcome into a stable code without
// inspecting provider causes.

type Cancelled struct{}

func (Cancelled) Error() string {
	return "the sign-in was cancelled"
}

func (Cancelled) reason() code.Code {
	return code.Cancelled
}

type Rejected struct{}

func (Rejected) Error() string {
	return "the sign-in could not be verified"
}

func (Rejected) reason() code.Code {
	return code.Rejected
}

type Unavailable struct{}

func (Unavailable) Error() string {
	return "authentication is temporarily unavailable"
}

func (Unavailable) reason() code.Code {
	return code.Unavailable
}

type Conflict struct{}

func (Conflict) Error() string {
	return "another account is signed in"
}

func (Conflict) reason() code.Code {
	return code.Conflict
}

type Storage struct{}

func (Storage) Error() string {
	return "secure credential storage is unavailable"
}

func (Storage) reason() code.Code {
	return code.Storage
}

type Forbidden struct{}

func (Forbidden) Error() string {
	return "the account has no active Drizz organization"
}

func (Forbidden) reason() code.Code {
	return code.Forbidden
}
