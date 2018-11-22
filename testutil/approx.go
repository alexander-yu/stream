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
