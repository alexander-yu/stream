package test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var precision = 9

func roundFloat(x float64, n int) float64 {
	unit := 5 * math.Pow10(-n-1)
	return math.Round(x/unit) * unit
}

// Approx asserts that two floats are approximately equal to each other,
// within 9 decimal points of precision.
func Approx(t *testing.T, x float64, y float64, msgAndArgs ...interface{}) {
	x = roundFloat(x, precision)
	y = roundFloat(y, precision)
	assert.Equal(t, x, y, msgAndArgs...)
}

// ApproxSlice asserts that two slices of floats are approximately equal to each other
// by element, within 9 decimal points of precision.
func ApproxSlice(t *testing.T, xs []float64, ys []float64) {
	require.Equal(t, len(xs), len(ys), "%v != %v: slice lengths are different", xs, ys)
	for i := 0; i < len(xs); i++ {
		Approx(t, xs[i], ys[i], "elements differ at index %d", i)
	}
}
