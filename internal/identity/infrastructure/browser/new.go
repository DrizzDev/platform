package browser

func New(options Options) (Browser, error) {
	if failure := options.validate(); failure != nil {
		return Browser{}, failure
	}
	return Browser{
		opener:  options.Opener,
		address: options.Address,
		path:    options.Path,
	}, nil
}
