package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestKurtosis(t *testing.T) {
	kurtosis, err := NewKurtosis()
	require.NoError(t, err)

	stream.TestData(kurtosis)

	value, err := kurtosis.Value()
	require.NoError(t, err)

	moment := 98. / 3.
	variance := 14. / 3.

	testutil.Approx(t, moment/math.Pow(variance, 2.)-3., value)
}
