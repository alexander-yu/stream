package joint

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/util/testutil"
)

func TestNewCorrelation(t *testing.T) {
	t.Run("pass: returns a Correlation", func(t *testing.T) {
		_, err := NewCorrelation(3)
		require.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewCorrelation(-1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestCorrelation(t *testing.T) {
	t.Run("pass: returns the correlation", func(t *testing.T) {
		correlation, err := NewCorrelation(3)
		require.NoError(t, err)

		err = testData(correlation)
		require.NoError(t, err)

		value, err := correlation.Value()
		require.NoError(t, err)

		testutil.Approx(t, 158./math.Sqrt(14.*5378./3.), value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		correlation, err := NewCorrelation(3)
		require.NoError(t, err)

		_, err = correlation.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: error if wrong number of values are pushed", func(t *testing.T) {
		correlation, err := NewCorrelation(3)
		require.NoError(t, err)

		vals := []float64{3.}
		err = correlation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Correlation expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))

		vals = []float64{3., 9., 27.}
		err = correlation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"Correlation expected 2 arguments: got %d (%v)",
			len(vals),
			vals,
		))
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		correlation, err := NewCorrelation(3)
		require.NoError(t, err)

		err = testData(correlation)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		correlation.core.queue.Dispose()
		err = correlation.Push(3., 9.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		correlation, err := NewCorrelation(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		correlation.core.queue.Dispose()
		vals := []float64{3., 9.}
		err = correlation.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %v to queue", vals))
	})
}
