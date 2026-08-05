package session

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	handle "github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

// Publication drives the two durable transitions of a renewal. Attempt marks
// the active revision one time before the network refresh, so a lost or
// concurrent renewal can never reuse the same refresh token. Publish writes the
// rotated candidate and advances the head under a fenced compare-and-swap.
type Publication interface {
	Attempt(context.Context, marking.Marking) error
	Publish(context.Context, Candidate) (Receipt, error)
}

// Candidate is the rotated credential to publish: the prior record it supersedes,
// the renewed secret, the epoch fenced by Attempt, and the moment of renewal.
type Candidate struct {
	Prior    credential.Record
	Renewal  Renewal
	Expected epoch.Epoch
	Moment   time.Time
}

// Receipt is the trusted context of the renewed credential.
type Receipt struct {
	Subject subject.Subject
	Session handle.Session
	Method  method.Method
	Expiry  time.Time
}
