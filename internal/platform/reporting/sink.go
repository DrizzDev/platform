package reporting

import (
	"context"
	"log/slog"
)

type Sink interface {
	Handler() slog.Handler
	Close(context.Context) error
}
