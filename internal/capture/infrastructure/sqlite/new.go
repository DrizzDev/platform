package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"go.opentelemetry.io/otel/metric"

	_ "modernc.org/sqlite"
)

func New(scope context.Context, options Options) (Store, error) {

	if failure := options.validate(); failure != nil {
		return Store{}, failure
	}
	latency, failure := options.Meter.Float64Histogram(instrument, metric.WithUnit("s"))

	if failure != nil {
		return Store{}, failure
	}

	settings := url.Values{}
	settings.Add("_txlock", "immediate")
	settings.Add("_pragma", "busy_timeout(2000)")
	settings.Add("_pragma", "foreign_keys(on)")
	settings.Add("_pragma", "journal_mode(wal)")
	settings.Add("_pragma", "synchronous(full)")
	settings.Add("_pragma", "max_page_count(2048)")

	handle, failure := sql.Open("sqlite", options.Path+"?"+settings.Encode())
	if failure != nil {
		return Store{}, failure
	}

	handle.SetMaxOpenConns(1)
	handle.SetMaxIdleConns(1)
	handle.SetConnMaxLifetime(0)

	store := Store{handle: handle, logger: options.Logger, tracer: options.Tracer, latency: latency}
	if failure := store.migrate(scope); failure != nil {
		return Store{}, errors.Join(failure, handle.Close())
	}
	return store, nil
}
