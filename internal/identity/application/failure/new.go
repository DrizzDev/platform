package failure

func New(input Input) (Value, error) {
	value := Value{
		code:        input.Code,
		retry:       input.Retry,
		detail:      input.Detail,
		correlation: input.Correlation,
	}
	if failure := value.validate(); failure != nil {
		return Value{}, failure
	}
	return value, nil
}
