package sentry

import (
	"errors"
	"math"
	"strconv"
)

const fallback = 1.0

func (input Input) parse() (float64, error) {
	if input.Sample == "" {
		return fallback, nil
	}
	sample, failure := strconv.ParseFloat(input.Sample, 64)
	if failure != nil || math.IsNaN(sample) || math.IsInf(sample, 0) || sample <= 0 || sample > 1 {
		return 0, errors.New("DRIZZ_SENTRY_SAMPLE_RATE must be a finite number greater than 0 and at most 1")
	}
	return sample, nil
}
