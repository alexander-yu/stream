package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestNewMean(t *testing.T) {
	t.Run("pass: returns a Mean", func(t *testing.T) {
		_, err := NewMean(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMean(-1)
		assert.EqualError(t, err, "-1 is a negative window")
	})
}

func TestMeanValue(t *testing.T) {
	t.Run("pass: returns the mean", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		stream.TestData(mean)

		value, err := mean.Value()
		require.NoError(t, err)

		testutil.Approx(t, 5, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		_, err = stream.SetupMetric(mean)
		require.NoError(t, err)

		_, err = mean.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
