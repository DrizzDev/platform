package reconcile

import (
	"context"
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

const backoff = time.Minute

// Reconciler drains the cleanup backlog a bounded batch at a time: each due
// candidate is deleted from the vault and acknowledged, or rescheduled when the
// vault is unavailable. It is best-effort and idempotent, so a partial pass is
// safe to repeat.
type Reconciler struct {
	queue Queue
	vault Vault
	clock Clock
}

func (reconciler Reconciler) Run(scope context.Context, _ Input) Result {
	records, failure := reconciler.queue.Pending(scope, reconciler.clock.Now())
	if failure != nil {
		return reconciler.deny(failure)
	}
	result := Result{}
	for _, record := range records {
		if reconciler.reclaim(scope, record) {
			result.reclaimed++
			continue
		}
		result.deferred++
	}
	return result
}

// reclaim deletes one candidate and acknowledges it, rescheduling the record
// when the vault is unavailable. The vault delete is idempotent, so a failed
// acknowledgement simply reclaims the same key on a later pass.
func (reconciler Reconciler) reclaim(scope context.Context, record cleanup.Record) bool {
	if failure := reconciler.vault.Delete(scope, credential.Key(record.Key())); failure != nil {
		reconciler.postpone(scope, record.Key())
		return false
	}
	_ = reconciler.queue.Acknowledge(scope, record.Key())
	return true
}

func (reconciler Reconciler) postpone(scope context.Context, key string) {
	entry, failure := deferral.New(deferral.Input{Key: key, Next: reconciler.clock.Now().Add(backoff)})
	if failure != nil {
		return
	}
	_ = reconciler.queue.Defer(scope, entry)
}

func (reconciler Reconciler) deny(cause error) Result {
	kind := code.Failed
	if errors.Is(cause, context.Canceled) || errors.Is(cause, context.DeadlineExceeded) {
		kind = code.Cancelled
	}
	value, _ := fault.New(fault.Input{Code: kind})
	return Result{failure: &value}
}
