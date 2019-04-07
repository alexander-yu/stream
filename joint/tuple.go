package joint

import (
	"math"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	mathutil "github.com/alexander-yu/stream/util/math"
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

func (m Tuple) hash() uint64 {
	result := 0
	// for practical purposes, the chance of collision is effectively
	// 0, since any joint moments are extremely unlikely to have individual
	// exponential terms that are higher than 4. As long as the individual elements
	// of the Tuple are less than 31, then it is impossible for collisions to occur
	// with this algorithm.
	for i := range m {
		result = 31*result + m[i]
	}
	return uint64(result)
}

func (m Tuple) eq(n Tuple) bool {
	if len(m) != len(n) {
		return false
	}

	for i := range m {
		if m[i] != n[i] {
			return false
		}
	}

	return true
}

func sub(m, n Tuple) (Tuple, error) {
	if len(m) != len(n) {
		return nil, errors.Errorf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		)
	}

	result := Tuple(make([]int, len(m)))
	for i := range m {
		result[i] = m[i] - n[i]
	}
	return result, nil
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

// iter executes a callback function over each Tuple that is less than
// or equal than the provided Tuple. For two Tuples m and n, we define
// m <= n iff m_i <= n_i for all i. The execution order is made by fixing
// the last element of the tuple first, and then incrementing the others
// until those options are exhausted. For example, for tuple = Tuple{2, 3},
// this is equivalent to the following (when iterating in increasing order):
//  for j := 0; j <= tuple[1], j++ {
//	    for i := 0; i <= tuple[0]; i++ {
//          cb(i, j)
//      }
//  }
// This execution order (rather than the expected one of having i on the
// outer loop) is due to the recursive nature of iter, and fact that it is
// faster to append arguments at the end rather than insert them at the beginning.
func iter(tuple Tuple, reverse bool, cb func(...int) error) error {
	if len(tuple) == 0 {
		return cb()
	}

	var result *multierror.Error
	if reverse {
		for i := tuple[len(tuple)-1]; i >= 0; i-- {
			err := iter(tuple[:len(tuple)-1], reverse, func(xs ...int) error {
				xs = append(xs, i)
				return cb(xs...)
			})
			result = multierror.Append(result, err)
		}
	} else {
		for i := 0; i <= tuple[len(tuple)-1]; i++ {
			err := iter(tuple[:len(tuple)-1], reverse, func(xs ...int) error {
				xs = append(xs, i)
				return cb(xs...)
			})
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
