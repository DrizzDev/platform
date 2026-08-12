package connect

import "log/slog"

type Options struct {
	Resolver resolver
	Store    store
	Recorder recorder
	Monitor  monitor
	Logger   *slog.Logger
}
