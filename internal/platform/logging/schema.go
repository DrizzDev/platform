package logging

import "log/slog"

type schema struct{}

func (schema schema) replace(input replacement) slog.Attr {
	if len(input.groups) != 0 {
		return input.attribute
	}
	switch input.attribute.Key {
	case slog.TimeKey:
		input.attribute.Key = timestamp
		input.attribute.Value = slog.TimeValue(input.attribute.Value.Time().UTC())
	case slog.MessageKey:
		input.attribute.Key = message
	}
	return input.attribute
}
