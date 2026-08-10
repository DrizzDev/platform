package operator

func New(options Options) (Operator, error) {
	if failure := options.validate(); failure != nil {
		return Operator{}, failure
	}
	return Operator{flow: options.Flow, recorder: options.Recorder, logger: options.Logger}, nil
}
