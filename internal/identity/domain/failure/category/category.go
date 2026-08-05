package category

type Category string

const (
	Account        Category = "ACCOUNT"
	Storage        Category = "STORAGE"
	Logout         Category = "LOGOUT"
	Internal       Category = "INTERNAL"
	Authorization  Category = "AUTHORIZATION"
	Authentication Category = "AUTHENTICATION"
)

func (category Category) Valid() bool {
	switch category {
	case Authentication, Authorization, Account, Storage, Logout, Internal:
		return true
	default:
		return false
	}
}

func (category Category) String() string {
	return string(category)
}
