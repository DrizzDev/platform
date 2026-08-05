package session

import "github.com/DrizzDev/platform/internal/identity/domain/failure/code"

// The port errors below are the session-owned failure contract. Infrastructure
// returns them so the flow can classify an outcome into a stable code without
// inspecting provider or storage causes.

// Missing reports that nothing is signed in.
type Missing struct{}

func (Missing) Error() string {
	return "no active credential is present"
}

func (Missing) reason() code.Code {
	return code.Required
}

// Uncertain reports that the credential can no longer be trusted or renewed: a
// fenced or superseded revision, or a refresh that cannot be safely retried. The
// only recovery is a new sign-in.
type Uncertain struct{}

func (Uncertain) Error() string {
	return "the session could not be renewed"
}

func (Uncertain) reason() code.Code {
	return code.Required
}

// Storage reports that the local credential store is unavailable.
type Storage struct{}

func (Storage) Error() string {
	return "secure credential storage is unavailable"
}

func (Storage) reason() code.Code {
	return code.Storage
}
