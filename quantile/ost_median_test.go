package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewOSTMedian(t *testing.T) {
	t.Run("pass: nonnegative window is valid", func(t *testing.T) {
		median, err := NewOSTMedian(0, AVL)
		require.NoError(t, err)
		assert.Equal(t, 0, median.quantile.window)

		median, err = NewOSTMedian(5, AVL)
		require.NoError(t, err)
		assert.Equal(t, 5, median.quantile.window)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewOSTMedian(-1, AVL)
		testutil.ContainsError(t, err, "error creating OSTQuantile")
	})

	t.Run("fail: unsupported OST implementation is invalid", func(t *testing.T) {
		_, err := NewOSTMedian(3, Impl(-1))
		testutil.ContainsError(t, err, "error creating OSTQuantile")
	})
}

func TestOSTMedianString(t *testing.T) {
	expectedString := fmt.Sprintf(
		"quantile.OSTMedian_{quantile:quantile.OSTQuantile_{window:3,interpolation:%d}}",
		Midpoint,
	)
	median, err := NewOSTMedian(3, AVL)
	require.NoError(t, err)

	assert.Equal(t, expectedString, median.String())
}

func TestOSTMedianPush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		median, err := NewOSTMedian(3, AVL)
		require.NoError(t, err)
		for i := 0.; i < 5; i++ {
			err := median.Push(i)
			require.NoError(t, err)
		}
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		median, err := NewOSTMedian(3, AVL)
		require.NoError(t, err)

		for i := 0.; i < 3; i++ {
			err = median.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		median.quantile.queue.Dispose()
		err = median.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		median, err := NewOSTMedian(3, AVL)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		median.quantile.queue.Dispose()
		val := 3.
		err = median.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestOSTMedianValue(t *testing.T) {
	t.Run("pass: if number of values is even, return average of middle two", func(t *testing.T) {
		median, err := NewOSTMedian(4, AVL)
		require.NoError(t, err)
		for i := 0.; i < 6; i++ {
			err := median.Push(i)
			require.NoError(t, err)
		}

		value, err := median.Value()
		require.NoError(t, err)

		assert.Equal(t, 3.5, value)
	})

	t.Run("pass: if number of values is odd, return middle value", func(t *testing.T) {
		median, err := NewOSTMedian(3, AVL)
		require.NoError(t, err)
		for i := 0.; i < 5; i++ {
			err := median.Push(i)
			require.NoError(t, err)
		}

		value, err := median.Value()
		require.NoError(t, err)

		assert.Equal(t, float64(3), value)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		median, err := NewOSTMedian(3, AVL)
		require.NoError(t, err)

		_, err = median.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})
}

func TestOSTMedianClear(t *testing.T) {
	median, err := NewOSTMedian(3, AVL)
	require.NoError(t, err)

	for i := 0.; i < 10; i++ {
		err = median.Push(i * i)
		require.NoError(t, err)
	}

	median.Clear()
	assert.Equal(t, uint64(0), median.quantile.queue.Len())
	assert.Equal(t, 0, median.quantile.statistic.Size())
}
