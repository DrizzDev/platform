package logging

import "log/slog"

type replacement struct {
	groups    []string
	attribute slog.Attr
}
