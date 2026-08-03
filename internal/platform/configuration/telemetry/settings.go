package telemetry

type Settings struct {
	exporter Exporter
	endpoint string
}

func (settings Settings) Exporter() Exporter {
	if settings.exporter == "" {
		return None
	}
	return settings.exporter
}

func (settings Settings) Endpoint() string {
	return settings.endpoint
}
