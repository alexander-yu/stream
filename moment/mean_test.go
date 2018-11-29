package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestMean(t *testing.T) {
	t.Run("pass: returns the mean", func(t *testing.T) {
		mean := &Mean{}

		stream.TestData(mean)

		value, err := mean.Value()
		require.NoError(t, err)

		testutil.Approx(t, 5, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		mean := &Mean{}

		stream.NewCore(&stream.CoreConfig{}, mean)

		_, err := mean.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
