package identifier

func New(value string) (Identifier, error) {
	identifier := Identifier{value: value}
	if failure := identifier.validate(); failure != nil {
		return Identifier{}, failure
	}
	return identifier, nil
}
