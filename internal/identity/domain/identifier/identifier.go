package identifier

type Identifier struct {
	value string
}

func (identifier Identifier) String() string {
	return identifier.value
}
