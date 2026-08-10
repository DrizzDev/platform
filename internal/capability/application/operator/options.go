package operator

import "log/slog"

type Options struct {
	Flow     flow
	Recorder recorder
	Logger   *slog.Logger
}
