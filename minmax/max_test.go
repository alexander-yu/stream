package minmax

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewMax(t *testing.T) {
	t.Run("pass: returns a Max", func(t *testing.T) {
		_, err := NewMax(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMax(-1)
		assert.EqualError(t, err, "-1 is a negative window")
	})
}

func TestMaxString(t *testing.T) {
	expectedString := "minmax.Max_{window:3}"
	max, err := NewMax(3)
	require.NoError(t, err)

	assert.Equal(t, expectedString, max.String())
}

func TestMaxValue(t *testing.T) {
	t.Run("pass: returns global maximum for a window of 0", func(t *testing.T) {
		max, err := NewMax(0)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 9., val)
	})

	t.Run("pass: returns maximum for a provided window", func(t *testing.T) {
		max, err := NewMax(5)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 5., val)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		max, err := NewMax(3)
		require.NoError(t, err)

		_, err = max.Value()
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		max, err := NewMax(3)
		require.NoError(t, err)

		for i := 0.; i < 3; i++ {
			err = max.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		max.queue.Dispose()
		err = max.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		max, err := NewMax(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		max.queue.Dispose()
		val := 3.
		err = max.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}
