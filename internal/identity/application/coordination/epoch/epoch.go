package epoch

type Epoch uint64

func (epoch Epoch) Next() Epoch {
	return epoch + 1
}

func (epoch Epoch) Before(other Epoch) bool {
	return epoch < other
}

func (epoch Epoch) After(other Epoch) bool {
	return epoch > other
}
