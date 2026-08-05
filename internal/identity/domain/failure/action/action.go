package action

type Action string

const (
	None   Action = "NONE"
	Retry  Action = "RETRY"
	Logout Action = "LOGOUT"
	Signin Action = "SIGN_IN"
)

func (action Action) Valid() bool {
	switch action {
	case Signin, Retry, Logout, None:
		return true
	default:
		return false
	}
}

func (action Action) String() string {
	return string(action)
}
