package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewMoment(t *testing.T) {
	moment := NewMoment(2, 3)
	assert.Equal(t, 2, moment.k)
	assert.Equal(t, 3, moment.window)
}

func TestMomentValue(t *testing.T) {
	t.Run("pass: returns the kth moment", func(t *testing.T) {
		moment := NewMoment(2, 3)
		err := Init(moment)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		value, err := moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 7, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		moment := NewMoment(2, 3)
		err := Init(moment)
		require.NoError(t, err)

		_, err = moment.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		moment := NewMoment(1, 3)
		err := Init(moment)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		moment.core.queue.Dispose()
		err = moment.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		moment := NewMoment(1, 3)
		err := Init(moment)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		moment.core.queue.Dispose()
		val := 3.
		err = moment.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		moment := NewMoment(1, 3)
		err := Init(moment)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		moment.Clear()
		expectedSums := []float64{0, 0}
		assert.Equal(t, expectedSums, moment.core.sums)
		assert.Equal(t, int(0), moment.core.count)
		assert.Equal(t, uint64(0), moment.core.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		moment := NewMoment(2, 3)
		expectedString := "moment.Moment_{k:2,window:3}"
		assert.Equal(t, expectedString, moment.String())
	})
}
