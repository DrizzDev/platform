package logging

import "log/slog"

type Policy struct{}

func (Policy) Handler() func([]string, slog.Attr) slog.Attr {
	return func(groups []string, attribute slog.Attr) slog.Attr {
		input := replacement{groups: groups, attribute: attribute}
		input.attribute = schema{}.replace(input)
		return redactor{}.replace(input)
	}
}
