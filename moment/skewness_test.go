package moment

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewSkewness(t *testing.T) {
	window := 3
	skewness := NewSkewness(window)
	assert.Equal(t, &Skewness{
		variance: &Moment{K: 2, Window: window},
		moment3:  &Moment{K: 3, Window: window},
		config: &CoreConfig{
			Sums: SumsConfig{
				2: true,
				3: true,
			},
			Window: &window,
		},
	}, skewness)
}

func TestSkewnessValue(t *testing.T) {
	t.Run("pass: returns the skewness", func(t *testing.T) {
		skewness := NewSkewness(3)
		err := SetupMetric(skewness)
		require.NoError(t, err)

		err = testData(skewness)
		require.NoError(t, err)

		value, err := skewness.Value()
		require.NoError(t, err)

		adjust := 3.
		moment := 9.
		variance := 7.

		testutil.Approx(t, adjust*moment/math.Pow(variance, 1.5), value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		skewness := NewSkewness(3)
		err := SetupMetric(skewness)
		require.NoError(t, err)

		_, err = skewness.Value()
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		skewness := NewSkewness(3)
		err := SetupMetric(skewness)
		require.NoError(t, err)

		err = testData(skewness)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		skewness.core.queue.Dispose()
		err = skewness.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		skewness := NewSkewness(3)
		err := SetupMetric(skewness)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		skewness.core.queue.Dispose()
		val := 3.
		err = skewness.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		skewness := NewSkewness(3)
		err := SetupMetric(skewness)
		require.NoError(t, err)

		err = testData(skewness)
		require.NoError(t, err)

		skewness.Clear()
		expectedSums := []float64{0, 0, 0, 0}
		assert.Equal(t, expectedSums, skewness.core.sums)
		assert.Equal(t, int(0), skewness.core.count)
		assert.Equal(t, uint64(0), skewness.core.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		expectedString := "moment.Skewness_{window:3}"
		skewness := NewSkewness(3)
		assert.Equal(t, expectedString, skewness.String())
	})
}
