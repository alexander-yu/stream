package moment

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewStd(t *testing.T) {
	t.Run("pass: returns an Std", func(t *testing.T) {
		_, err := NewStd(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewStd(-1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestStdValue(t *testing.T) {
	t.Run("pass: returns the standard deviation", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		err = testData(std)
		require.NoError(t, err)

		value, err := std.Value()
		require.NoError(t, err)

		testutil.Approx(t, math.Sqrt(7.), value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		_, err = std.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		err = testData(std)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		std.variance.core.queue.Dispose()
		err = std.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		std.variance.core.queue.Dispose()
		val := 3.
		err = std.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		err = testData(std)
		require.NoError(t, err)

		std.Clear()
		expectedSums := []float64{0, 0, 0}
		assert.Equal(t, expectedSums, std.variance.core.sums)
		assert.Equal(t, int(0), std.variance.core.count)
		assert.Equal(t, uint64(0), std.variance.core.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		expectedString := "moment.Std_{window:3}"
		std, err := NewStd(3)
		require.NoError(t, err)

		assert.Equal(t, expectedString, std.String())
	})
}
