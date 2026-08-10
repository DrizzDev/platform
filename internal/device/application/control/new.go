package control

func New(options Options) (Flow, error) {
	if failure := options.validate(); failure != nil {
		return Flow{}, failure
	}
	return Flow{bridge: options.Bridge}, nil
}
