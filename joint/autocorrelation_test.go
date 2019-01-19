package joint

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewAutocorrelation(t *testing.T) {
	autocorrelation := NewAutocorrelation(1, 3)
	assert.Equal(t, 1, autocorrelation.lag)
	assert.Equal(t, NewCorrelation(3), autocorrelation.correlation)
}

func TestAutocorrelation(t *testing.T) {
	t.Run("pass: returns the autocorrelation", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		value, err := autocorrelation.Value()
		require.NoError(t, err)

		testutil.Approx(t, 31.*math.Sqrt(2289.)/1526., value)
	})

	t.Run("pass: returns the correlation for lag=0", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(0, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		value, err := autocorrelation.Value()
		require.NoError(t, err)

		testutil.Approx(t, 158./math.Sqrt(14.*5378./3.), value)
	})

	t.Run("fail: if core push fails for lag=0, return error", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(0, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to push to the core
		autocorrelation.core.queue.Dispose()
		err = autocorrelation.Push(3., 9.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error popping item from queue",
			3.,
			9.,
		))
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		_, err = autocorrelation.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: error if wrong number of values are pushed", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		vals := []float64{3.}
		err = autocorrelation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Autocorrelation expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))

		vals = []float64{3., 9., 27.}
		err = autocorrelation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Autocorrelation expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))
	})

	t.Run("fail: if core queue retrieval fails, return error", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		autocorrelation.core.queue.Dispose()
		err = autocorrelation.Push(3., 9.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error popping item from queue",
			3.,
			64.,
		))
	})

	t.Run("fail: if core queue insertion fails, return error", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		autocorrelation.core.queue.Dispose()

		// no error yet because we have not populated the lag yet
		err = autocorrelation.Push(8., 64.)
		require.NoError(t, err)

		err = autocorrelation.Push(3., 9.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error pushing %v to queue",
			3.,
			64.,
			[]float64{3., 64.},
		))
	})

	t.Run("fail: if lag queue retrieval fails, return error", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		autocorrelation.queue.Dispose()
		err = autocorrelation.Push(3., 9.)
		testutil.ContainsError(t, err, "error popping item from lag queue")
	})

	t.Run("fail: if lag queue insertion fails, return error", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		autocorrelation.queue.Dispose()
		err = autocorrelation.Push(8., 64.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing %f to lag queue",
			64.,
		))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		err := Init(autocorrelation)
		require.NoError(t, err)

		err = testData(autocorrelation)
		require.NoError(t, err)

		autocorrelation.Clear()

		expectedSums := map[uint64]float64{
			0:  0.,
			1:  0.,
			2:  0.,
			31: 0.,
			32: 0.,
			62: 0.,
		}
		assert.Equal(t, expectedSums, autocorrelation.core.sums)
		assert.Equal(t, expectedSums, autocorrelation.core.newSums)
		assert.Equal(t, 0, autocorrelation.core.count)
		assert.Equal(t, uint64(0), autocorrelation.core.queue.Len())
		assert.Equal(t, uint64(0), autocorrelation.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		autocorrelation := NewAutocorrelation(1, 3)
		expectedString := "joint.Autocorrelation_{lag:1,window:3}"
		assert.Equal(t, expectedString, autocorrelation.String())
	})
}
