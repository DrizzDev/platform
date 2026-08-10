package recording

func New(options Options) (Recorder, error) {
	if failure := options.validate(); failure != nil {
		return Recorder{}, failure
	}
	return Recorder{
		sink:   options.Sink,
		clock:  options.Clock,
		lease:  options.Lease,
		logger: options.Logger,
		writer: options.Writer,
		keeper: options.Keeper,
	}, nil
}
