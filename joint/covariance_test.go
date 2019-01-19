package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestCovarianceValue(t *testing.T) {
	t.Run("pass: returns the covariance", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		err = testData(covariance)
		require.NoError(t, err)

		value, err := covariance.Value()
		require.NoError(t, err)

		testutil.Approx(t, 79., value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		_, err = covariance.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: error if wrong number of values are pushed", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		vals := []float64{3.}
		err = covariance.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Covariance expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))

		vals = []float64{3., 9., 27.}
		err = covariance.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Covariance expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		err = testData(covariance)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		covariance.core.queue.Dispose()
		err = covariance.Push(3., 9.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		covariance.core.queue.Dispose()
		vals := []float64{3., 9.}
		err = covariance.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %v to queue", vals))
	})

	t.Run("pass: Clear() resets the metric", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		err := Init(covariance)
		require.NoError(t, err)

		err = testData(covariance)
		require.NoError(t, err)

		covariance.Clear()

		expectedSums := map[uint64]float64{
			0:  0.,
			1:  0.,
			31: 0.,
			32: 0.,
		}
		assert.Equal(t, expectedSums, covariance.core.sums)
		assert.Equal(t, expectedSums, covariance.core.newSums)
		assert.Equal(t, 0, covariance.core.count)
		assert.Equal(t, uint64(0), covariance.core.queue.Len())
	})

	t.Run("pass: String() returns string representation", func(t *testing.T) {
		covariance := &Covariance{Window: 3}
		expectedString := "joint.Covariance_{window:3}"
		assert.Equal(t, expectedString, covariance.String())
	})
}
