package logging

import (
	"log/slog"

	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
)

func New(options Options) (*slog.Logger, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	return slog.New(options.handler()).With(
		slog.String(identity, options.Build.Name()),
		slog.String(version, options.Build.Version()),
		slog.String(revision, options.Build.Revision()),
	), nil
}

// handler emits structured JSON at the configured level, or discards everything when logging is off, so the output
// destination stays fixed and the configured level alone decides visibility.
func (options Options) handler() slog.Handler {
	if options.Settings.Level() == logging.Off {
		return slog.DiscardHandler
	}
	return slog.NewJSONHandler(options.Output, &slog.HandlerOptions{
		AddSource:   true,
		Level:       options.level(),
		ReplaceAttr: Policy{}.Handler(),
	})
}
