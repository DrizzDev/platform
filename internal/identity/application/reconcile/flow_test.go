package reconcile_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/reconcile"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

type queue struct {
	records      []cleanup.Record
	fault        error
	acknowledged []string
	deferred     []string
}

func (queue *queue) Pending(context.Context, time.Time) ([]cleanup.Record, error) {
	return queue.records, queue.fault
}

func (queue *queue) Acknowledge(_ context.Context, key string) error {
	queue.acknowledged = append(queue.acknowledged, key)
	return nil
}

func (queue *queue) Defer(_ context.Context, entry deferral.Deferral) error {
	queue.deferred = append(queue.deferred, entry.Key())
	return nil
}

type locker struct {
	stuck map[string]bool
}

func (locker locker) Delete(_ context.Context, key credential.Key) error {
	if locker.stuck[key.String()] {
		return errors.New("vault locked")
	}
	return nil
}

type clock struct{}

func (clock) Now() time.Time { return time.Unix(5000, 0) }

type harness struct {
	test *testing.T
}

func (harness harness) record(key string) cleanup.Record {
	harness.test.Helper()
	record, failure := cleanup.New(cleanup.Input{
		Key: key, Reason: cleanup.Superseded, State: cleanup.Pending,
		Next: time.Unix(1000, 0), Deadline: time.Unix(9000, 0), Created: time.Unix(1000, 0),
	})
	if failure != nil {
		harness.test.Fatal(failure)
	}
	return record
}

func (harness harness) build(options reconcile.Options) reconcile.Reconciler {
	harness.test.Helper()
	made, failure := reconcile.New(options)
	if failure != nil {
		harness.test.Fatal(failure)
	}
	return made
}

func TestReclaim(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	store := &queue{records: []cleanup.Record{fixture.record("LOCAL#1#a"), fixture.record("LOCAL#2#b")}}
	result := fixture.build(reconcile.Options{Queue: store, Vault: locker{}, Clock: clock{}}).
		Run(context.Background(), reconcile.Input{})

	if result.Failed() || result.Reclaimed() != 2 || result.Deferred() != 0 {
		test.Fatalf("result = %+v", result)
	}
	if len(store.acknowledged) != 2 || len(store.deferred) != 0 {
		test.Fatalf("acknowledged %v deferred %v", store.acknowledged, store.deferred)
	}
}

func TestDefer(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	store := &queue{records: []cleanup.Record{fixture.record("LOCAL#1#a"), fixture.record("LOCAL#2#b")}}
	result := fixture.build(reconcile.Options{
		Queue: store, Vault: locker{stuck: map[string]bool{"LOCAL#2#b": true}}, Clock: clock{},
	}).Run(context.Background(), reconcile.Input{})

	if result.Reclaimed() != 1 || result.Deferred() != 1 {
		test.Fatalf("result = %+v", result)
	}
	if len(store.acknowledged) != 1 || store.acknowledged[0] != "LOCAL#1#a" {
		test.Fatalf("acknowledged = %v", store.acknowledged)
	}
	if len(store.deferred) != 1 || store.deferred[0] != "LOCAL#2#b" {
		test.Fatalf("deferred = %v", store.deferred)
	}
}

func TestBacklogFault(test *testing.T) {
	test.Parallel()

	fixture := harness{test: test}
	store := &queue{fault: errors.New("database is locked")}
	result := fixture.build(reconcile.Options{Queue: store, Vault: locker{}, Clock: clock{}}).
		Run(context.Background(), reconcile.Input{})

	fault, present := result.Failure()
	if !present || fault.Code() != code.Failed {
		test.Fatalf("code = %v (present %v)", fault.Code(), present)
	}
}
