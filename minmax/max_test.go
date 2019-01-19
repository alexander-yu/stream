package minmax

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewMax(t *testing.T) {
	max := NewMax(3)
	assert.Equal(t, 3, max.window)
	assert.Equal(t, uint64(0), max.queue.Len())
	assert.Equal(t, 0, max.deque.Len())
	assert.Equal(t, math.Inf(-1), max.max)
	assert.Equal(t, 0, max.count)
}

func TestMaxString(t *testing.T) {
	expectedString := "minmax.Max_{window:3}"
	max := NewMax(3)
	assert.Equal(t, expectedString, max.String())
}

func TestMaxValue(t *testing.T) {
	t.Run("pass: returns global maximum for a window of 0", func(t *testing.T) {
		max := NewMax(0)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err := max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 9., val)
	})

	t.Run("pass: returns maximum for a provided window", func(t *testing.T) {
		max := NewMax(5)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err := max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 5., val)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		max := NewMax(3)
		_, err := max.Value()
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		max := NewMax(3)

		for i := 0.; i < 3; i++ {
			err := max.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		max.queue.Dispose()
		err := max.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		max := NewMax(3)

		// dispose the queue to simulate an error when we try to insert into the queue
		max.queue.Dispose()
		val := 3.
		err := max.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestMaxClear(t *testing.T) {
	max := NewMax(3)

	for i := 0.; i < 3; i++ {
		err := max.Push(i)
		require.NoError(t, err)
	}

	max.Clear()
	assert.Equal(t, 0, max.count)
	assert.Equal(t, math.Inf(-1), max.max)
	assert.Equal(t, uint64(0), max.queue.Len())
	assert.Equal(t, 0, max.deque.Len())
}
