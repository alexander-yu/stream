package joint

import (
	"fmt"
	"math"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func testAutocorrData(autocorr *Autocorrelation) error {
	for i := 1.; i < 5; i++ {
		err := autocorr.Push(i)
		if err != nil {
			return errors.Wrapf(err, "failed to push %f to metric", i)
		}
	}

	err := autocorr.Push(8.)
	if err != nil {
		return errors.Wrapf(err, "failed to push %f to metric", 8.)
	}
	return nil
}

func TestNewAutocorrelation(t *testing.T) {
	t.Run("pass: valid Autocorrelation is valid", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)
		assert.Equal(t, 1, autocorrelation.lag)
		assert.Equal(t, NewCorrelation(3), autocorrelation.correlation)
	})

	t.Run("fail: negative lag returns error", func(t *testing.T) {
		_, err := NewAutocorrelation(-1, 3)
		testutil.ContainsError(t, err, "-1 is a negative lag")
	})
}

func TestAutocorrelation(t *testing.T) {
	t.Run("pass: returns the autocorrelation", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
		require.NoError(t, err)

		value, err := autocorrelation.Value()
		require.NoError(t, err)

		testutil.Approx(t, 5./(2.*math.Sqrt(7.)), value)
	})

	t.Run("pass: returns the correlation for lag=0", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(0, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
		require.NoError(t, err)

		value, err := autocorrelation.Value()
		require.NoError(t, err)

		testutil.Approx(t, 1., value)
	})

	t.Run("fail: if core push fails for lag=0, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(0, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to push to the core
		autocorrelation.core.queue.Dispose()
		err = autocorrelation.Push(3.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error popping item from queue",
			3.,
			3.,
		))
	})

	t.Run("fail: if Core is not set, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = autocorrelation.Push(0.)
		testutil.ContainsError(t, err, "Core is not set")

		_, err = autocorrelation.Value()
		testutil.ContainsError(t, err, "Core is not set")
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		_, err = autocorrelation.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: error if wrong number of values are pushed", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		vals := []float64{3., 5.}
		err = autocorrelation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Autocorrelation expected 1 argument: got %d (%v)",
			len(vals),
			vals,
		))

		vals = []float64{}
		err = autocorrelation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Autocorrelation expected 1 argument: got %d (%v)",
			len(vals),
			vals,
		))
	})

	t.Run("fail: if core queue retrieval fails, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		autocorrelation.core.queue.Dispose()
		err = autocorrelation.Push(3.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error popping item from queue",
			3.,
			8.,
		))
	})

	t.Run("fail: if core queue insertion fails, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		autocorrelation.core.queue.Dispose()

		// no error yet because we have not populated the lag yet
		err = autocorrelation.Push(8.)
		require.NoError(t, err)

		err = autocorrelation.Push(3.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing (%f, %f) to core: error pushing %v to queue",
			3.,
			8.,
			[]float64{3., 8.},
		))
	})

	t.Run("fail: if lag queue retrieval fails, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		autocorrelation.queue.Dispose()
		err = autocorrelation.Push(3.)
		testutil.ContainsError(t, err, "error popping item from lag queue")
	})

	t.Run("fail: if lag queue insertion fails, return error", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		autocorrelation.queue.Dispose()
		err = autocorrelation.Push(8.)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"error pushing %f to lag queue",
			8.,
		))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		err = Init(autocorrelation)
		require.NoError(t, err)

		err = testAutocorrData(autocorrelation)
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
		autocorrelation, err := NewAutocorrelation(1, 3)
		require.NoError(t, err)

		expectedString := "joint.Autocorrelation_{lag:1,window:3}"
		assert.Equal(t, expectedString, autocorrelation.String())
	})
}
