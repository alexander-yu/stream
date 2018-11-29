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
	_, err := NewKurtosis()
	assert.NoError(t, err)
}

func TestKurtosisValue(t *testing.T) {
	t.Run("pass: returns the excess kurtosis", func(t *testing.T) {
		kurtosis, err := NewKurtosis()
		require.NoError(t, err)

		stream.TestData(kurtosis)

		value, err := kurtosis.Value()
		require.NoError(t, err)

		moment := 98. / 3.
		variance := 14. / 3.

		testutil.Approx(t, moment/math.Pow(variance, 2.)-3., value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		kurtosis, err := NewKurtosis()
		require.NoError(t, err)

		stream.NewCore(&stream.CoreConfig{}, kurtosis)

		_, err = kurtosis.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
