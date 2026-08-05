package logout

import "github.com/DrizzDev/platform/internal/identity/domain/failure/code"

// Missing reports that nothing is signed in, which the flow treats as an
// already-complete logout rather than a failure.
type Missing struct{}

func (Missing) Error() string {
	return "no active credential is present"
}

// Storage reports that the local credential store is unavailable.
type Storage struct{}

func (Storage) Error() string {
	return "secure credential storage is unavailable"
}

func (Storage) reason() code.Code {
	return code.Storage
}
