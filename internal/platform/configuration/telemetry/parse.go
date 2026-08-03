package telemetry

import "strings"

func New(input Input) (Settings, error) {
	exporter := Exporter(strings.ToUpper(input.Exporter))
	if exporter == "" {
		exporter = None
	}
	settings := Settings{exporter: exporter, endpoint: input.Endpoint}
	if failure := settings.validate(); failure != nil {
		return Settings{}, failure
	}
	return settings, nil
}
