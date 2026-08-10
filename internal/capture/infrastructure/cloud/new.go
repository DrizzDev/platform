package cloud

func New(options Options) (Client, error) {
	if failure := options.validate(); failure != nil {
		return Client{}, failure
	}
	return Client{
		base:   options.Base,
		agent:  options.Agent,
		source: &source{provider: options.Provider, clock: options.Clock},
	}, nil
}
