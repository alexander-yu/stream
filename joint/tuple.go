package joint

import (
	"math"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream/util/mathutil"
)

// Tuple represents a vector of integers, which
// for the purpose of this package represents a
// vector of exponents for multivariate moments.
type Tuple []int

func (m Tuple) abs() int {
	sum := 0
	for _, i := range m {
		sum += i
	}
	return sum
}

func (m Tuple) hash() int {
	result := 0
	// for practical purposes, the chance of collision is effectively
	// 0, since any joint moments are extremely unlikely to have individual
	// exponential terms that are higher than 4. As long as the individual elements
	// of the Tuple are less than 31, then it is impossible for collisions to occur
	// with this algorithm.
	for i := range m {
		result = 31*result + m[i]
	}
	return result
}

func multinom(m, n Tuple) (int, error) {
	if len(m) != len(n) {
		return 0, errors.Errorf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		)
	}

	result := 1
	for i := range m {
		result *= mathutil.Binom(m[i], n[i])
	}

	return result, nil
}

func pow(x []float64, n Tuple) (float64, error) {
	if len(x) != len(n) {
		return 0, errors.Errorf(
			"Cannot exponentiate slice and Tuple with different lengths: %d != %d",
			len(x),
			len(n),
		)
	}

	result := 1.
	for i := range x {
		result *= math.Pow(x[i], float64(n[i]))
	}

	return result, nil
}
