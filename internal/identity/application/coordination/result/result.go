package result

type Result string

const (
	Published Result = "PUBLISHED"
	Rejected  Result = "REJECTED"
	Uncertain Result = "UNCERTAIN"
)

func (result Result) Valid() bool {
	switch result {
	case Published, Rejected, Uncertain:
		return true
	default:
		return false
	}
}

func (result Result) String() string {
	return string(result)
}
