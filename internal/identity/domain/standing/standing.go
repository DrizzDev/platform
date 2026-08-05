package standing

type Standing string

const (
	Active  Standing = "ACTIVE"
	Expired Standing = "EXPIRED"
	Revoked Standing = "REVOKED"
)

func (standing Standing) Valid() bool {
	switch standing {
	case Active, Expired, Revoked:
		return true
	default:
		return false
	}
}

func (standing Standing) String() string {
	return string(standing)
}
