package courier

import (
	"bytes"
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Run drains every pending entry in order, uploading each and flipping it to synced. It stops at the first failure,
// leaving the remaining entries pending and recoverable on the next run.
func (courier Courier) Run(scope context.Context) error {
	pending, failure := courier.ledger.Pending(scope)

	if failure != nil {
		return failure
	}
	for _, entry := range pending {
		if failure := courier.ship(scope, entry); failure != nil {
			return failure
		}
	}
	return nil
}

// ship uploads one entry's artifact, then the entry, then records it synced. The artifact goes first so the entry
// never references bytes the cloud has not received.
func (courier Courier) ship(scope context.Context, entry journal.Entry) error {
	if reference := entry.Artifact(); !reference.Empty() {
		blob, failure := courier.vault.Get(scope, reference)
		if failure != nil {
			return failure
		}
		failure = courier.retry(scope, func(scope context.Context) error {
			return courier.uploader.Blob(scope, Cargo{Digest: reference, Source: bytes.NewReader(blob)})
		})
		if failure != nil {
			return failure
		}
	}
	failure := courier.retry(scope, func(scope context.Context) error {
		return courier.uploader.Record(scope, entry)
	})
	if failure != nil {
		return failure
	}
	return courier.ledger.Advance(scope, journal.Transition{Entry: entry, State: journal.Synced})
}
