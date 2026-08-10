package recording

import "log/slog"

type Options struct {
	Sink   sink
	Writer writer
	Logger *slog.Logger
}
