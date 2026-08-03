package logging

import (
	"log/slog"

	"github.com/DrizzDev/platform/internal/platform/configuration/logging"
)

func (options Options) level() slog.Level {
	switch options.Settings.Level() {

	case logging.Debug:
		return slog.LevelDebug

	case logging.Info:
		return slog.LevelInfo

	case logging.Warn:
		return slog.LevelWarn

	case logging.Error:
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}
