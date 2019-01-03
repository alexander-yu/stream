package quantile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/quantile/heap"
	"github.com/alexander-yu/stream/util/testutil"
)

func TestHeapMedianPush(t *testing.T) {
	median := NewHeapMedian()
	for i := 0.; i < 5; i++ {
		err := median.Push(i)
		require.NoError(t, err)
	}

	testutil.ApproxSlice(t, []float64{1., 0.}, median.lowHeap.Values())
	testutil.ApproxSlice(t, []float64{2., 3., 4.}, median.highHeap.Values())

	err := median.Push(3.)
	require.NoError(t, err)

	testutil.ApproxSlice(t, []float64{2., 0., 1.}, median.lowHeap.Values())
	testutil.ApproxSlice(t, []float64{3., 3., 4.}, median.highHeap.Values())

	err = median.Push(2.)
	require.NoError(t, err)

	testutil.ApproxSlice(t, []float64{2., 2., 1., 0.}, median.lowHeap.Values())
	testutil.ApproxSlice(t, []float64{3., 3., 4.}, median.highHeap.Values())

	err = median.Push(1.)
	require.NoError(t, err)

	testutil.ApproxSlice(t, []float64{2., 1., 1., 0.}, median.lowHeap.Values())
	testutil.ApproxSlice(t, []float64{2., 3., 4., 3.}, median.highHeap.Values())
}

func TestHeapMedianValue(t *testing.T) {
	t.Run("pass: if low heap is larger, return its top", func(t *testing.T) {
		median := NewHeapMedian()
		median.lowHeap = heap.NewHeap([]float64{2., 2., 1., 0.}, fmax)
		median.highHeap = heap.NewHeap([]float64{3., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2., value)
	})

	t.Run("pass: if high heap is larger, return its top", func(t *testing.T) {
		median := NewHeapMedian()
		median.lowHeap = heap.NewHeap([]float64{1., 0.}, fmax)
		median.highHeap = heap.NewHeap([]float64{2., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2., value)
	})

	t.Run("pass: if heaps are equal in size, return average of tops", func(t *testing.T) {
		median := NewHeapMedian()
		median.lowHeap = heap.NewHeap([]float64{2., 0., 1.}, fmax)
		median.highHeap = heap.NewHeap([]float64{3., 3., 4.}, fmin)

		value, err := median.Value()
		require.NoError(t, err)

		testutil.Approx(t, 2.5, value)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		median := NewHeapMedian()
		_, err := median.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}