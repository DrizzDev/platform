package console

func New(options Options) (Display, error) {
	if failure := options.validate(); failure != nil {
		return Display{}, failure
	}
	return Display{writer: options.Writer}, nil
}
