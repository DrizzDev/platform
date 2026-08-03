package logging

import (
	"errors"
	"strings"
)

func New(value string) (Settings, error) {
	switch Level(strings.ToUpper(value)) {
	case "":
		return Settings{level: Info}, nil
	case Debug:
		return Settings{level: Debug}, nil
	case Info:
		return Settings{level: Info}, nil
	case Warn:
		return Settings{level: Warn}, nil
	case Error:
		return Settings{level: Error}, nil
	default:
		return Settings{}, errors.New("DRIZZ_LOG_LEVEL must be one of debug, info, warn, error")
	}
}
