package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrument = "drizz.identity.duration"
	boundary   = "identity.store"
	operation  = "drizz.operation"
	outcome    = "drizz.outcome"
)

type verdict string

const (
	settled  verdict = "OK"
	declined verdict = "DECLINED"
	faulted  verdict = "FAULT"
)

type unit = func(context.Context, *sql.Tx) error

type reader = func(context.Context) error

type task struct {
	work unit
	name string
}

type probe struct {
	work reader
	name string
}

// transact runs a write as one short immediate transaction, observed end to end.
// A returned error rolls the transaction back; a nil return commits it.
func (store Store) transact(scope context.Context, task task) error {
	return store.observe(scope, probe{name: task.name, work: func(scope context.Context) error {
		transaction, failure := store.handle.BeginTx(scope, nil)
		if failure != nil {
			return failure
		}
		defer func() { _ = transaction.Rollback() }()
		if failure := task.work(scope, transaction); failure != nil {
			return failure
		}
		return transaction.Commit()
	}})
}

// observe wraps an operation with a span and a duration metric. A fault records
// only the operation code, never the driver cause, which stays with the returned
// error for the application boundary to map.
func (store Store) observe(scope context.Context, probe probe) (failure error) {
	scope, span := store.tracer.Start(scope, boundary+"."+probe.name)
	start := time.Now()
	defer func() {
		result := store.grade(failure)
		if result == faulted {
			span.SetStatus(codes.Error, probe.name)
			store.logger.ErrorContext(scope, "identity.store.failed", slog.String(operation, probe.name))
		}
		store.latency.Record(scope, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String(operation, probe.name),
			attribute.String(outcome, string(result))))
		span.End()
	}()
	failure = probe.work(scope)
	return failure
}

func (Store) grade(failure error) verdict {
	var expected benign
	switch {
	case failure == nil:
		return settled
	case errors.As(failure, &expected):
		return declined
	default:
		return faulted
	}
}
