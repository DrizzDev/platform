package agent

// Base names the well-known directory an agent's configuration path is anchored to. The domain names the anchor
// without resolving it, so no home-directory lookup — an operating-system call — happens in the core; an infrastructure
// adapter turns a Base into a real directory for the running user.
type Base string

const (
	// Home is the user's home directory, e.g. ~/.claude.json or ~/.codex/config.toml.
	Home Base = "HOME"
	// Config is the user's per-application configuration directory, e.g. macOS "Library/Application Support" or
	// Linux ~/.config, under which desktop agents keep their settings.
	Config Base = "CONFIG"
)

func (base Base) Valid() bool {
	switch base {
	case Home, Config:
		return true
	default:
		return false
	}
}

func (base Base) String() string {
	return string(base)
}
