package configuration

const (
	prefix = "DRIZZ_"

	// reserved is the device-helper override namespace read by the device layer for tests and continuous integration.
	// It is not a configuration setting, so keys under it are deliberately not treated as unknown settings.
	reserved = "DRIZZ_DEVICE_"

	level    = "DRIZZ_LOG_LEVEL"
	exporter = "DRIZZ_TELEMETRY_EXPORTER"
	endpoint = "DRIZZ_TELEMETRY_ENDPOINT"
	dsn      = "DRIZZ_SENTRY_DSN"
	sample   = "DRIZZ_SENTRY_SAMPLE_RATE"
	stage    = "DRIZZ_SENTRY_ENVIRONMENT"

	issuer   = "DRIZZ_AUTH0_ISSUER"
	client   = "DRIZZ_AUTH0_CLIENT"
	audience = "DRIZZ_AUTH0_AUDIENCE"
	redirect = "DRIZZ_AUTH0_REDIRECT"
	scopes   = "DRIZZ_AUTH0_SCOPES"
	session  = "DRIZZ_SESSION"
	cloud    = "DRIZZ_CLOUD"
)
