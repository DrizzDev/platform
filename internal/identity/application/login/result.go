package login

import (
	"time"

	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/standing"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

// Result is the trusted context of a completed login, or a code-only failure.
// It never carries a token, provider response, or raw cause.
type Result struct {
	failure  *fault.Value
	subject  subject.Subject
	session  session.Session
	method   method.Method
	standing standing.Standing
	expiry   time.Time
	tenant   Tenant
}

func (result Result) Organization() string {
	return result.tenant.Name
}

func (result Result) Subject() subject.Subject {
	return result.subject
}

func (result Result) Session() session.Session {
	return result.session
}

func (result Result) Method() method.Method {
	return result.method
}

func (result Result) Standing() standing.Standing {
	return result.standing
}

func (result Result) Expiry() time.Time {
	return result.expiry
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
