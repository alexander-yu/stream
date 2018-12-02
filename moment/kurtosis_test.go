package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestNewKurtosis(t *testing.T) {
	t.Run("pass: returns a Kurtosis", func(t *testing.T) {
		_, err := NewKurtosis(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewKurtosis(-1)
		assert.EqualError(t, err, "error creating 2nd Moment: -1 is a negative window")
	})
}

func TestKurtosisValue(t *testing.T) {
	t.Run("pass: returns the excess kurtosis", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		stream.TestData(kurtosis)

		value, err := kurtosis.Value()
		require.NoError(t, err)

		moment := 98. / 3.
		variance := 14. / 3.

		testutil.Approx(t, moment/math.Pow(variance, 2.)-3., value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		_, err = stream.SetupMetric(kurtosis)
		require.NoError(t, err)

		_, err = kurtosis.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
