package publication

import (
	"context"
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/retraction"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/logout"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

var (
	_ logout.Vault       = Publisher{}
	_ logout.Publication = Publisher{}
)

// Active reads the credential the active pointer names, reporting Missing when
// nothing is signed in so logout stays idempotent.
func (publisher Publisher) Active(scope context.Context) (credential.Record, error) {
	head, failure := publisher.ledger.Head(scope, publisher.session)
	var absent sqlite.Absent
	switch {
	case errors.As(failure, &absent):
		return credential.Record{}, logout.Missing{}
	case failure != nil:
		return credential.Record{}, logout.Storage{}
	}
	record, failure := publisher.vault.Read(scope, credential.Key(head.Key()))
	if failure != nil {
		return credential.Record{}, logout.Storage{}
	}
	return record, nil
}

// Retract clears the active pointer and queues its key for deletion. A pointer
// already gone is treated as a completed retraction.
func (publisher Publisher) Retract(scope context.Context, moment time.Time) error {
	request, failure := retraction.New(retraction.Input{Session: publisher.session, Moment: moment})
	if failure != nil {
		return failure
	}
	failure = publisher.ledger.Retract(scope, request)
	var absent sqlite.Absent
	if errors.As(failure, &absent) {
		return nil
	}
	return failure
}
