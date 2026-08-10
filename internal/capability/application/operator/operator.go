package operator

import "log/slog"

// Operator performs one catalogued capability: it drives the device through the neutral port and records the
// execution, so the command line and the agent connection produce the same authoritative record from the same code.
type Operator struct {
	flow     flow
	recorder recorder
	logger   *slog.Logger
}
