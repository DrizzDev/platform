package operator

import "log/slog"

// Operator performs one catalogued capability: it drives the device through the neutral port, records the execution,
// and instruments it, so the command line and the agent connection produce the same authoritative record and the same
// telemetry from the same code.
type Operator struct {
	flow     flow
	recorder recorder
	logger   *slog.Logger
	monitor  monitor
}
