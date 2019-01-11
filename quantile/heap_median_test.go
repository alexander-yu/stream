package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/quantile/heap"
	"github.com/alexander-yu/stream/util/testutil"
)

func TestNewHeapMedian(t *testing.T) {
	t.Run("pass: returns a HeapMedian", func(t *testing.T) {
		_, err := NewHeapMedian(10)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewHeapMedian(-1)
		assert.EqualError(t, err, "-1 is a negative window")
	})
}

func TestHeapMedianString(t *testing.T) {
	expectedString := "quantile.HeapMedian_{window:3}"
	median, err := NewHeapMedian(3)
	require.NoError(t, err)

	assert.Equal(t, expectedString, median.String())
}

func TestHeapMedianPush(t *testing.T) {
	t.Run("pass: maintains heaps properly", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		for i := 0.; i < 3; i++ {
			err = median.Push(i)
			require.NoError(t, err)

			err = median.Push(10 - i)
			require.NoError(t, err)
		}

		testutil.ApproxSlice(t, []float64{2, 0, 1}, median.lowHeap.Values())
		testutil.ApproxSlice(t, []float64{8, 10, 9}, median.highHeap.Values())

		err = median.Push(1)
		require.NoError(t, err)

		err = median.Push(1)
		require.NoError(t, err)

		err = median.Push(1)
		require.NoError(t, err)

		testutil.ApproxSlice(t, []float64{1, 1, 1, 0, 1}, median.lowHeap.Values())
		testutil.ApproxSlice(t, []float64{2, 8, 9, 10}, median.highHeap.Values())

		for i := 0.; i < 3; i++ {
			err = median.Push(i)
			require.NoError(t, err)

			err = median.Push(10 - i)
			require.NoError(t, err)
		}

		testutil.ApproxSlice(t, []float64{1, 1, 1, 0, 1}, median.lowHeap.Values())
		testutil.ApproxSlice(t, []float64{2, 8, 9, 10, 8}, median.highHeap.Values())

		err = median.Push(9)
		require.NoError(t, err)

		err = median.Push(1)
		require.NoError(t, err)

		testutil.ApproxSlice(t, []float64{1, 1, 1, 0, 1}, median.lowHeap.Values())
		testutil.ApproxSlice(t, []float64{2, 8, 9, 10, 9}, median.highHeap.Values())
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = median.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		median.queue.Dispose()
		err = median.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		median.queue.Dispose()
		val := 3.
		err = median.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})

	t.Run("fail: if queue insertion fails (after queue retrieval), return error", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = median.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to insert into the queue
		median.queue.Dispose()
		err = median.Push(3.) // TODO: get this test case to work
		testutil.ContainsError(t, err, "error popping item from queue")
	})
}

func TestHeapMedianValue(t *testing.T) {
	t.Run("pass: if low heap is larger, return its top", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		median.lowHeap = heap.NewHeap("low", []float64{2., 2., 1., 0.}, fmax)
		median.highHeap = heap.NewHeap("high", []float64{3., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2., value)
	})

	t.Run("pass: if high heap is larger, return its top", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		median.lowHeap = heap.NewHeap("low", []float64{1., 0.}, fmax)
		median.highHeap = heap.NewHeap("high", []float64{2., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2., value)
	})

	t.Run("pass: if heaps are equal in size, return average of tops", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		median.lowHeap = heap.NewHeap("low", []float64{2., 0., 1.}, fmax)
		median.highHeap = heap.NewHeap("high", []float64{3., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2.5, value)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		median, err := NewHeapMedian(10)
		require.NoError(t, err)

		_, err = median.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
