package recording

func New(options Options) (Recorder, error) {
	if failure := options.validate(); failure != nil {
		return Recorder{}, failure
	}
	return Recorder{writer: options.Writer, sink: options.Sink, logger: options.Logger}, nil
}
