package release

func New(input Input) (Identity, error) {
	identity := Identity{
		name:     input.Name,
		version:  input.Version,
		revision: input.Revision,
	}
	if failure := identity.validate(); failure != nil {
		return Identity{}, failure
	}
	return identity, nil
}
