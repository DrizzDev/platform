package login

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

// Publication is the login-owned port that persists the validated credential. It
// resolves the session identity and revision, writes the vault candidate, and
// runs the fenced compare-and-swap. It reports Conflict when a different account
// is already active and Storage when the vault is unavailable.
type Publication interface {
	Publish(context.Context, Candidate) (Receipt, error)
	Retract(context.Context, time.Time) error
}

// Candidate is the validated provider credential to persist. The session
// identity and revision are resolved by the adapter under the fence.
type Candidate struct {
	Token  Token
	Moment time.Time
}

// Receipt is the trusted context of the published credential.
type Receipt struct {
	Subject subject.Subject
	Session session.Session
	Method  method.Method
	Expiry  time.Time
}
