package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestNewMoment(t *testing.T) {
	t.Run("pass: returns a Kurtosis", func(t *testing.T) {
		moment, err := NewMoment(0, 3)
		require.NoError(t, err)
		assert.Equal(t, 0, moment.k)

		moment, err = NewMoment(5, 3)
		require.NoError(t, err)
		assert.Equal(t, 5, moment.k)
	})

	t.Run("fail: negative moment is invalid", func(t *testing.T) {
		_, err := NewMoment(-1, 3)
		assert.EqualError(t, err, "-1 is a negative moment")
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMoment(3, -1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestMoment(t *testing.T) {
	t.Run("pass: returns the kth moment", func(t *testing.T) {
		moment, err := NewMoment(2, 3)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		value, err := moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 7, value)
	})

	t.Run("pass: 0th moment always returns 1", func(t *testing.T) {
		moment, err := NewMoment(0, 3)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		value, err := moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 1, value)

		err = moment.core.Push(10)
		require.NoError(t, err)

		value, err = moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 1, value)
	})

	t.Run("pass: 1st moment always returns 0", func(t *testing.T) {
		moment, err := NewMoment(1, 3)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		value, err := moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 0, value)

		err = moment.core.Push(10)
		require.NoError(t, err)

		value, err = moment.Value()
		require.NoError(t, err)

		testutil.Approx(t, 0, value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		moment, err := NewMoment(2, 3)
		require.NoError(t, err)

		_, err = moment.Value()
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		moment, err := NewMoment(1, 3)
		require.NoError(t, err)

		err = testData(moment)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		moment.core.queue.Dispose()
		err = moment.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		moment, err := NewMoment(1, 3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		moment.core.queue.Dispose()
		val := 3.
		err = moment.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})
}
