package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewOSTQuantile(t *testing.T) {
	t.Run("fail: unsupported OST implementation is invalid", func(t *testing.T) {
		impl := Impl(-1)
		config := &Config{
			Window:        stream.IntPtr(3),
			Interpolation: Linear.Ptr(),
			Impl:          &impl,
		}
		_, err := NewOSTQuantile(config)
		testutil.ContainsError(t, err, "error validating config")
	})
}

func TestOSTQuantileString(t *testing.T) {
	expectedString := fmt.Sprintf(
		"quantile.OSTQuantile_{window:3,interpolation:%d}",
		Linear,
	)
	config := &Config{
		Window:        stream.IntPtr(3),
		Interpolation: Linear.Ptr(),
		Impl:          AVL.Ptr(),
	}
	quantile, err := NewOSTQuantile(config)
	require.NoError(t, err)

	assert.Equal(t, expectedString, quantile.String())
}

func TestOSTQuantilePush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(3),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(3),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(3),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		quantile.queue.Dispose()
		val := 3.
		err = quantile.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestOSTQuantileValue(t *testing.T) {
	t.Run("pass: returns quantile for exact index", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(5),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Lower.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Higher.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Nearest.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Nearest.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(8),
			Interpolation: Nearest.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Nearest.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Midpoint.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
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
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
		require.NoError(t, err)

		_, err = quantile.Value(0.25)
		testutil.ContainsError(t, err, "no values seen yet")
	})

	t.Run("fail: if quantile not in (0, 1), return error", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(6),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		quantile, err := NewOSTQuantile(config)
		require.NoError(t, err)

		_, err = quantile.Value(0.)
		testutil.ContainsError(t, err, fmt.Sprintf("quantile %f not in (0, 1)", 0.))

		_, err = quantile.Value(1.)
		testutil.ContainsError(t, err, fmt.Sprintf("quantile %f not in (0, 1)", 1.))
	})
}

func TestOSTQuantileClear(t *testing.T) {
	config := &Config{
		Window:        stream.IntPtr(3),
		Interpolation: Linear.Ptr(),
		Impl:          AVL.Ptr(),
	}
	quantile, err := NewOSTQuantile(config)
	require.NoError(t, err)

	for i := 0.; i < 10; i++ {
		err = quantile.Push(i * i)
		require.NoError(t, err)
	}

	quantile.Clear()
	assert.Equal(t, uint64(0), quantile.queue.Len())
	assert.Equal(t, 0, quantile.statistic.Size())
}
