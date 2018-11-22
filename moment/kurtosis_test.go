package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"stream"
)

func TestKurtosis(t *testing.T) {
	kurtosis, err := NewKurtosis()
	require.NoError(err)

	core := stream.TestData(kurtosis)
	core.Push(8)

	value, err := kurtosis.Value()
	require.NoError(t, err)

	moment := 98. / 3.
	variance := 14. / 3.

	approx(t, moment/math.Pow(variance, 2.)-3., value)
}
