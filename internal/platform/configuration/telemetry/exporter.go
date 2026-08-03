package telemetry

type Exporter string

const (
	None Exporter = "NONE"
	OTLP Exporter = "OTLP"
)
