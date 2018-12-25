package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/util/testutil"
)

func TestNewMean(t *testing.T) {
	t.Run("pass: returns a Mean", func(t *testing.T) {
		_, err := NewMean(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMean(-1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestMeanValue(t *testing.T) {
	t.Run("pass: returns the mean", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		err = testData(mean)
		require.NoError(t, err)

		value, err := mean.Value()
		require.NoError(t, err)

		testutil.Approx(t, 5, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		_, err = mean.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		err = testData(mean)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		mean.core.queue.Dispose()
		err = mean.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		mean, err := NewMean(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		mean.core.queue.Dispose()
		val := 3.
		err = mean.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})
}
