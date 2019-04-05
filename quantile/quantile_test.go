package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewQuantile(t *testing.T) {
	t.Run("fail: invalid Option is invalid", func(t *testing.T) {
		_, err := NewQuantile(3, ImplOption(-1))
		testutil.ContainsError(t, err, "error setting option")
	})
}

func TestQuantileString(t *testing.T) {
	expectedString := fmt.Sprintf(
		"quantile.Quantile_{window:3,interpolation:%d}",
		Linear,
	)
	quantile, err := NewQuantile(3)
	require.NoError(t, err)

	assert.Equal(t, expectedString, quantile.String())
}

func TestQuantilePush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		quantile, err := NewQuantile(3)
		require.NoError(t, err)
		for i := 0.; i < 5; i++ {
			err := quantile.Push(i)
			require.NoError(t, err)
		}

		assert.Equal(t, uint64(3), quantile.queue.Len())
		for i := 2.; i < 5; i++ {
			val, err := quantile.queue.Get()
			y := val.(float64)
			require.NoError(t, err)
			testutil.Approx(t, i, y)
		}

		assert.Equal(t, 3, quantile.statistic.Size())
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		quantile, err := NewQuantile(3)
		require.NoError(t, err)

		for i := 0.; i < 3; i++ {
			err = quantile.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		quantile.queue.Dispose()
		err = quantile.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		quantile, err := NewQuantile(3)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		quantile.queue.Dispose()
		val := 3.
		err = quantile.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestQuantileValue(t *testing.T) {
	t.Run("pass: returns quantile for exact index", func(t *testing.T) {
		quantile, err := NewQuantile(5)
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, 36., value)
	})

	t.Run("pass: returns quantile with linear interpolation", func(t *testing.T) {
		quantile, err := NewQuantile(6)
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, .75*25+.25*36, value)
	})

	t.Run("pass: returns quantile with lower interpolation", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Lower))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, 25., value)
	})

	t.Run("pass: returns quantile with higher interpolation", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Higher))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, 36., value)
	})

	t.Run("pass: returns quantile with nearest interpolation (delta < 0.5)", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Nearest))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, 25., value)
	})

	t.Run("pass: returns quantile with nearest interpolation (delta == 0.5, idx % 2 == 0)", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Nearest))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.5)
		require.NoError(t, err)
		testutil.Approx(t, 36., value)
	})

	t.Run("pass: returns quantile with nearest interpolation (delta == 0.5, idx % 2 == 1)", func(t *testing.T) {
		quantile, err := NewQuantile(8, InterpolationOption(Nearest))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.5)
		require.NoError(t, err)
		testutil.Approx(t, 36., value)
	})

	t.Run("pass: returns quantile with nearest interpolation (delta > 0.5)", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Nearest))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.75)
		require.NoError(t, err)
		testutil.Approx(t, 64., value)
	})

	t.Run("pass: returns quantile with midpoint interpolation", func(t *testing.T) {
		quantile, err := NewQuantile(6, InterpolationOption(Midpoint))
		require.NoError(t, err)

		for i := 0.; i < 10; i++ {
			err = quantile.Push(i * i)
			require.NoError(t, err)
		}

		value, err := quantile.Value(0.25)
		require.NoError(t, err)
		testutil.Approx(t, 30.5, value)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		quantile, err := NewQuantile(6)
		require.NoError(t, err)

		_, err = quantile.Value(0.25)
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if quantile not in (0, 1), return error", func(t *testing.T) {
		quantile, err := NewQuantile(6)
		require.NoError(t, err)

		_, err = quantile.Value(0.)
		testutil.ContainsError(t, err, fmt.Sprintf("quantile %f not in (0, 1)", 0.))

		_, err = quantile.Value(1.)
		testutil.ContainsError(t, err, fmt.Sprintf("quantile %f not in (0, 1)", 1.))
	})
}

func TestQuantileClear(t *testing.T) {
	quantile, err := NewQuantile(3)
	require.NoError(t, err)

	for i := 0.; i < 10; i++ {
		err = quantile.Push(i * i)
		require.NoError(t, err)
	}

	quantile.Clear()
	assert.Equal(t, uint64(0), quantile.queue.Len())
	assert.Equal(t, 0, quantile.statistic.Size())
}

func TestQuantileRLock(t *testing.T) {
	quantile, err := NewQuantile(3)
	require.NoError(t, err)

	done := make(chan bool)

	err = quantile.Push(1.)
	require.NoError(t, err)

	// Lock for reading
	quantile.RLock()

	// spawn a goroutine to push; should be blocked until RUnlock() is called
	go func() {
		err := quantile.Push(3.)
		require.NoError(t, err)
		done <- true
	}()

	val, err := quantile.Value(0.5)
	require.NoError(t, err)
	testutil.Approx(t, 1., val)

	// Undo RLock call
	quantile.RUnlock()

	// New Push call should now be unblocked
	<-done
	val, err = quantile.Value(0.5)
	require.NoError(t, err)
	testutil.Approx(t, 2., val)
}
