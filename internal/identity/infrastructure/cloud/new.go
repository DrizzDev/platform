package cloud

func New(options Options) (Client, error) {
	if failure := options.validate(); failure != nil {
		return Client{}, failure
	}
	client := Client{agent: options.Agent, base: options.Base}
	if options.Provider != nil {
		client.source = &source{provider: options.Provider, clock: options.Clock}
	}
	return client, nil
}
