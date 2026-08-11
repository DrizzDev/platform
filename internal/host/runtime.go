package host

import (
	"context"
	"fmt"

	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration"
	"github.com/DrizzDev/platform/internal/platform/observability"
	"github.com/DrizzDev/platform/internal/transport/mcp"
)

type runtime struct {
	environment []string
	streams     Streams
	identity    release.Identity
	build       build.Info
}

func (runtime runtime) Run(scope context.Context) error {
	settings, failure := configuration.New(runtime.environment).Load()
	if failure != nil {
		return runtime.reject(failure)
	}
	observer, failure := observability.New(scope, observability.Options{
		Build:    runtime.build,
		Settings: settings,
		Output:   runtime.streams.Failure,
	})
	if failure != nil {
		return runtime.fail(failure)
	}
	return runtime.serve(scope, observer)
}

func (runtime runtime) serve(scope context.Context, observer observability.Provider) error {
	current := session{observer: observer}
	defer current.shutdown(scope)
	server, failure := mcp.New(mcp.Options{
		Release:  runtime.identity,
		Logger:   observer.Diagnostics(),
		External: observer.External(),
		Tracer:   observer.Tracer(),
		Meter:    observer.Meter(),
		Input:    runtime.streams.Input,
		Output:   runtime.streams.Output,
		Perform: pilot{foundation{
			environment: runtime.environment,
			streams:     runtime.streams,
			build:       runtime.build,
		}},
	})
	if failure != nil {
		observer.Report(scope)
		return handled{cause: failure}
	}
	if failure := server.Run(scope); failure != nil {
		return handled{cause: failure}
	}
	return nil
}

// reject surfaces a configuration error. Configuration errors are Drizz-authored,
// actionable, and never echo secret values (CODE-027), so they are shown.
func (runtime runtime) reject(failure error) error {
	_, _ = fmt.Fprintln(runtime.streams.Failure, failure)
	return handled{cause: failure}
}

// fail surfaces a provider construction failure as a code. The cause may carry
// endpoint or machine detail, so it is retained only for control flow.
func (runtime runtime) fail(failure error) error {
	_, _ = fmt.Fprintln(runtime.streams.Failure, "startup.failed")
	return handled{cause: failure}
}
