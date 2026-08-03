package release

type Identity struct {
	name     string
	version  string
	revision string
}

func (identity Identity) Name() string {
	return identity.name
}

func (identity Identity) Version() string {
	return identity.version
}

func (identity Identity) Revision() string {
	return identity.revision
}
