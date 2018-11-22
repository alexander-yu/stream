package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"stream"
)

func TestSkewness(t *testing.T) {
	skewness, err := NewSkewness()
	require.NoError(t, err)

	core := stream.TestData(skewness)
	core.Push(8)

	value, err := skewness.Value()
	require.NoError(t, err)

	adjust := 3.
	moment := 9.
	variance := 7.

	approx(t, adjust*moment/math.Pow(variance, 1.5), value)
}
