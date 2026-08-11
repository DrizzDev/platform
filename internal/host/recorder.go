package host

import (
	"context"
	"path/filepath"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/system"
	"github.com/DrizzDev/platform/internal/platform/observability"
)

const (
	// lease is how long an in-progress recording is protected from routine cleanup while it is being written.
	lease = 5 * time.Minute

	// artifacts is the directory, beside the capture database, that holds the content-addressed captured images.
	artifacts = "artifacts"
)

// recorder assembles the local recorder from the capture stores — the ordered journal and the content-addressed
// artifact store, both under Drizz's per-user data directory. The returned store is closed by the caller on teardown.
func (foundation foundation) recorder(scope context.Context, observer observability.Provider) (recording.Recorder, sqlite.Store, error) {
	path, failure := sqlite.Location{}.Resolve()
	if failure != nil {
		return recording.Recorder{}, sqlite.Store{}, failure
	}
	store, failure := sqlite.New(scope, sqlite.Options{
		Path:   path,
		Logger: observer.Diagnostics(),
		Tracer: observer.Tracer(),
		Meter:  observer.Meter(),
	})
	if failure != nil {
		return recording.Recorder{}, sqlite.Store{}, failure
	}
	vault, failure := artifact.New(artifact.Options{
		Root:   filepath.Join(filepath.Dir(path), artifacts),
		Logger: observer.Diagnostics(),
		Tracer: observer.Tracer(),
		Meter:  observer.Meter(),
	})
	if failure != nil {
		_ = store.Close()
		return recording.Recorder{}, sqlite.Store{}, failure
	}
	made, failure := recording.New(recording.Options{
		Sink:   vault,
		Writer: store,
		Keeper: store,
		Clock:  system.New(),
		Logger: observer.Diagnostics(),
		Lease:  lease,
	})
	if failure != nil {
		_ = store.Close()
		return recording.Recorder{}, sqlite.Store{}, failure
	}
	return made, store, nil
}
