package agent

// Dialect is the on-disk format of an agent application's configuration file. The merge engine picks a reader and a
// writer by dialect, so a new agent that uses an existing format needs only a catalog row, not new code.
type Dialect string

const (
	Json Dialect = "JSON"
	Toml Dialect = "TOML"
)

func (dialect Dialect) Valid() bool {
	switch dialect {
	case Json, Toml:
		return true
	default:
		return false
	}
}

func (dialect Dialect) String() string {
	return string(dialect)
}
