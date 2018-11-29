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
	_, err := NewSkewness()
	assert.NoError(t, err)
}

func TestSkewness(t *testing.T) {
	t.Run("pass: returns the skewness", func(t *testing.T) {
		skewness, err := NewSkewness()
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
		skewness, err := NewSkewness()
		require.NoError(t, err)

		stream.NewCore(&stream.CoreConfig{}, skewness)

		_, err = skewness.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
