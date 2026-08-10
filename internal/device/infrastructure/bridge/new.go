package bridge

func New(options Options) (*Channel, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	channel := &Channel{
		options: options,
		done:    make(chan struct{}),
		slots:   make(chan struct{}, inflight),
	}
	channel.build = channel.spawn
	return channel, nil
}
