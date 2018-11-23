package testutil

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var precision = 9

func roundFloat(x float64, n int) float64 {
	unit := 5 * math.Pow10(-n-1)
	return math.Round(x/unit) * unit
}

// Approx asserts that two floats are approximately equal to each other,
// within 9 decimal points of precision.
func Approx(t *testing.T, x float64, y float64) {
	x = roundFloat(x, precision)
	y = roundFloat(y, precision)
	assert.Equal(t, x, y)
}

// ApproxSlice asserts that two slices of floats are approximately equal to each other
// by element, within 9 decimal points of precision.
func ApproxSlice(t *testing.T, xs []float64, ys []float64) {
	assert.Equal(t, len(xs), len(ys))

	// only check up to the smaller of the lists
	var num int
	if len(xs) < len(ys) {
		num = len(xs)
	} else {
		num = len(ys)
	}

	for i := 0; i < num; i++ {
		Approx(t, xs[i], ys[i])
	}
}
