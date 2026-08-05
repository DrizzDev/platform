package login

func New(options Options) (Flow, error) {
	if failure := options.validate(); failure != nil {
		return Flow{}, failure
	}
	return Flow{
		establishment: options.Establishment,
		publication:   options.Publication,
		authority:     options.Authority,
		clock:         options.Clock,
	}, nil
}
