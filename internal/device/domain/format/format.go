package format

// Format is the encoding of a captured frame. The supported set is the vocabulary the sidecar guarantees;
// new encodings extend it additively.
type Format string

const (
	Png Format = "PNG"
)

func (format Format) Valid() bool {
	switch format {
	case Png:
		return true
	default:
		return false
	}
}

func (format Format) String() string {
	return string(format)
}
