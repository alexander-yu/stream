package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestMeanValue(t *testing.T) {
	t.Run("pass: returns the mean", func(t *testing.T) {
		mean := &Mean{Window: 3}
		err := Init(mean)
		require.NoError(t, err)

		err = testData(mean)
		require.NoError(t, err)

		value, err := mean.Value()
		require.NoError(t, err)

		testutil.Approx(t, 5, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		mean := &Mean{Window: 3}
		err := Init(mean)
		require.NoError(t, err)

		_, err = mean.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		mean := &Mean{Window: 3}
		err := Init(mean)
		require.NoError(t, err)

		err = testData(mean)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		mean.core.queue.Dispose()
		err = mean.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		mean := &Mean{Window: 3}
		err := Init(mean)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		mean.core.queue.Dispose()
		val := 3.
		err = mean.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		mean := &Mean{Window: 3}
		err := Init(mean)
		require.NoError(t, err)

		err = testData(mean)
		require.NoError(t, err)

		mean.Clear()
		assert.Equal(t, float64(0), mean.core.mean)
		assert.Equal(t, uint64(0), mean.core.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		expectedString := "moment.Mean_{window:3}"
		mean := &Mean{Window: 3}

		assert.Equal(t, expectedString, mean.String())
	})
}
