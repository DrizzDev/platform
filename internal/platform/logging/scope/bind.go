package scope

import "log/slog"

func (scope Scope) Bind(logger *slog.Logger) *slog.Logger {
	if scope.request != "" {
		logger = logger.With(slog.String(request, scope.request))
	}
	if scope.correlation != "" {
		logger = logger.With(slog.String(correlation, scope.correlation))
	}
	if scope.session != "" {
		logger = logger.With(slog.String(session, scope.session))
	}
	return logger
}
