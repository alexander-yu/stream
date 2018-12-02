package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestNewSkewness(t *testing.T) {
	t.Run("pass: returns a Skewness", func(t *testing.T) {
		_, err := NewSkewness(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewSkewness(-1)
		assert.EqualError(t, err, "error creating 2nd Moment: -1 is a negative window")
	})
}

func TestSkewness(t *testing.T) {
	t.Run("pass: returns the skewness", func(t *testing.T) {
		skewness, err := NewSkewness(3)
		require.NoError(t, err)

		stream.TestData(skewness)

		value, err := skewness.Value()
		require.NoError(t, err)

		adjust := 3.
		moment := 9.
		variance := 7.

		testutil.Approx(t, adjust*moment/math.Pow(variance, 1.5), value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		skewness, err := NewSkewness(3)
		require.NoError(t, err)

		_, err = stream.SetupMetric(skewness)
		require.NoError(t, err)

		_, err = skewness.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
