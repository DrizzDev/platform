package logging

type Level string

const (
	Off   Level = "OFF"
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
)
