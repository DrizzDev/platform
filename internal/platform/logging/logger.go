package logging

import (
	"log/slog"
)

func New(options Options) (*slog.Logger, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	handler := slog.NewJSONHandler(options.Output, &slog.HandlerOptions{
		AddSource:   true,
		Level:       options.level(),
		ReplaceAttr: Policy{}.Handler(),
	})
	return slog.New(handler).With(
		slog.String(identity, options.Build.Name()),
		slog.String(version, options.Build.Version()),
		slog.String(revision, options.Build.Revision()),
	), nil
}
