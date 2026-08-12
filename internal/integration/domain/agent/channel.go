package agent

// Channel is how an agent hands a hook program the event: on standard input as JSON, or as a single JSON argument.
// The receiver reads the event the way the calling agent delivers it.
type Channel string

const (
	Stdin Channel = "STDIN"
	Argv  Channel = "ARGV"
)

func (channel Channel) String() string {
	return string(channel)
}
