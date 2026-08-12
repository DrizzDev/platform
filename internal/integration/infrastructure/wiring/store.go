// Package wiring reads and writes agent application configuration files. It resolves each agent's path from its
// descriptor, parses the file in its own dialect into one generic document, and merges the Drizz entry without
// disturbing any other setting, publishing every change atomically with a backup of the original.
package wiring

import (
	"context"
	"errors"
	"path/filepath"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// Store implements the installer's per-agent configuration port for every supported dialect. It carries a tracer and a
// latency metric so each configuration-file write is measured at the filesystem provider boundary.
type Store struct {
	tracer   trace.Tracer
	duration metric.Float64Histogram
}

type Options struct {
	Tracer trace.Tracer
	Meter  metric.Meter
}

func New(options Options) (Store, error) {
	if options.Tracer == nil || options.Meter == nil {
		return Store{}, errors.New("wiring tracer and meter are required")
	}
	duration, failure := options.Meter.Float64Histogram("drizz.wiring.duration", metric.WithUnit("s"))
	if failure != nil {
		return Store{}, failure
	}
	return Store{tracer: options.Tracer, duration: duration}, nil
}

// Detect reports whether the agent application is present on this machine: its configuration file exists, or — for an
// agent that keeps its file in a dedicated directory — that directory exists.
func (store Store) Detect(target agent.Agent) (bool, error) {
	path, failure := store.locate(target)
	if failure != nil {
		return false, failure
	}
	if store.present(path) {
		return true, nil
	}
	if len(target.Segments()) > 1 && store.present(filepath.Dir(path)) {
		return true, nil
	}
	return false, nil
}

// Wired reports whether Drizz is already registered in the agent's configuration.
func (store Store) Wired(target agent.Agent) (bool, error) {
	document, failure := store.read(target)
	if failure != nil {
		return false, failure
	}
	servers, _ := document[target.Collection()].(map[string]any)
	if servers == nil {
		return false, nil
	}
	_, found := servers[connect.Name]
	return found, nil
}

// Connect merges the Drizz entry into the agent's configuration, creating the file or the server collection if
// absent, preserving every other entry, and confirming the write by reading it back.
func (store Store) Connect(scope context.Context, job connect.Task) (failure error) {
	scope, gauge := store.begin(scope, "connect")
	defer func() { gauge.close(scope, failure) }()

	document, failure := store.read(job.Agent)
	if failure != nil {
		return failure
	}
	collection := job.Agent.Collection()
	servers, _ := document[collection].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}
	servers[job.Server.Name()] = store.entry(job)
	document[collection] = servers

	coder, failure := store.codec(job.Agent.Dialect())
	if failure != nil {
		return failure
	}
	raw, failure := coder.render(document)
	if failure != nil {
		return connect.Locked{}
	}
	path, failure := store.locate(job.Agent)
	if failure != nil {
		return failure
	}
	if failure := store.stamp(change{path: path, raw: raw}); failure != nil {
		return failure
	}
	return store.confirm(job)
}

// Disconnect removes only the Drizz entry from the agent's configuration, dropping the now-empty server collection but
// leaving every other setting untouched. Removing an entry that is not there succeeds.
func (store Store) Disconnect(scope context.Context, target agent.Agent) (failure error) {
	scope, gauge := store.begin(scope, "disconnect")
	defer func() { gauge.close(scope, failure) }()

	document, failure := store.read(target)
	if failure != nil {
		return failure
	}
	collection := target.Collection()
	servers, _ := document[collection].(map[string]any)
	if servers == nil {
		return nil
	}
	if _, found := servers[connect.Name]; !found {
		return nil
	}
	delete(servers, connect.Name)
	if len(servers) == 0 {
		delete(document, collection)
	} else {
		document[collection] = servers
	}

	coder, failure := store.codec(target.Dialect())
	if failure != nil {
		return failure
	}
	raw, failure := coder.render(document)
	if failure != nil {
		return connect.Locked{}
	}
	path, failure := store.locate(target)
	if failure != nil {
		return failure
	}
	return store.stamp(change{path: path, raw: raw})
}

func (store Store) confirm(job connect.Task) error {
	document, failure := store.read(job.Agent)
	if failure != nil {
		return failure
	}
	servers, _ := document[job.Agent.Collection()].(map[string]any)
	if servers == nil {
		return connect.Locked{}
	}
	if _, found := servers[job.Server.Name()]; !found {
		return connect.Locked{}
	}
	return nil
}
