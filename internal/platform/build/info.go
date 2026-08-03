package build

type Info struct {
	name     string
	version  string
	revision string
	modified bool
}

func (info Info) Name() string {
	return info.name
}

func (info Info) Version() string {
	return info.version
}

func (info Info) Revision() string {
	return info.revision
}

func (info Info) Modified() bool {
	return info.modified
}
