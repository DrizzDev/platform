package lobby

func New(options Options) (Lobby, error) {
	if failure := options.validate(); failure != nil {
		return Lobby{}, failure
	}
	return Lobby{
		clock:    options.Clock,
		window:   options.span(),
		capacity: options.size(),
		register: options.Register,
	}, nil
}
