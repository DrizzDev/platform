package reporting

import (
	"github.com/DrizzDev/platform/internal/platform/reporting/sentry"
)

func New(options Options) (Provider, error) {
	var sinks []Sink
	if options.Settings.Sentry().Enabled() {
		sink, failure := sentry.New(sentry.Options{
			Build:    options.Build,
			Settings: options.Settings.Sentry(),
		})
		if failure != nil {
			return Provider{}, failure
		}
		sinks = append(sinks, sink)
	}
	return Provider{sinks: sinks}, nil
}
