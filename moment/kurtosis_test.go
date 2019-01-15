package moment

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewKurtosis(t *testing.T) {
	t.Run("pass: returns a Kurtosis", func(t *testing.T) {
		_, err := NewKurtosis(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewKurtosis(-1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestKurtosisString(t *testing.T) {
	expectedString := "moment.Kurtosis_{window:3}"
	kurtosis, err := NewKurtosis(3)
	require.NoError(t, err)

	assert.Equal(t, expectedString, kurtosis.String())
}

func TestKurtosisValue(t *testing.T) {
	t.Run("pass: returns the excess kurtosis", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		err = testData(kurtosis)
		require.NoError(t, err)

		value, err := kurtosis.Value()
		require.NoError(t, err)

		moment := 98. / 3.
		variance := 14. / 3.

		testutil.Approx(t, moment/math.Pow(variance, 2.)-3., value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		_, err = kurtosis.Value()
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		err = testData(kurtosis)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		kurtosis.core.queue.Dispose()
		err = kurtosis.Push(3.)
		testutil.ContainsError(t, err, "error pushing to core: error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		kurtosis, err := NewKurtosis(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		kurtosis.core.queue.Dispose()
		val := 3.
		err = kurtosis.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing to core: error pushing %f to queue", val))
	})
}

func TestKurtosisClear(t *testing.T) {
	kurtosis, err := NewKurtosis(3)
	require.NoError(t, err)

	err = testData(kurtosis)
	require.NoError(t, err)

	kurtosis.Clear()
	expectedSums := []float64{0, 0, 0, 0, 0}
	assert.Equal(t, expectedSums, kurtosis.core.sums)
	assert.Equal(t, int(0), kurtosis.core.count)
	assert.Equal(t, uint64(0), kurtosis.core.queue.Len())
}
