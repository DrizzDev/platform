package release

import "errors"

func (identity Identity) validate() error {
	switch {
	case identity.name == "":
		return errors.New("release name is required")
	case identity.version == "":
		return errors.New("release version is required")
	case identity.revision == "":
		return errors.New("release revision is required")
	default:
		return nil
	}
}
