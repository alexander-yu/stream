package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewIQR(t *testing.T) {
	t.Run("pass: nonnegative window is valid", func(t *testing.T) {
		iqr, err := NewIQR(0, AVL)
		require.NoError(t, err)
		assert.Equal(t, 0, iqr.quantile.window)

		iqr, err = NewIQR(5, AVL)
		require.NoError(t, err)
		assert.Equal(t, 5, iqr.quantile.window)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewIQR(-1, AVL)
		testutil.ContainsError(t, err, "error creating Quantile")
	})

	t.Run("fail: unsupported Impl is invalid", func(t *testing.T) {
		_, err := NewIQR(3, Impl(-1))
		testutil.ContainsError(t, err, "error creating Quantile")
	})
}

func TestIQRString(t *testing.T) {
	expectedString := fmt.Sprintf(
		"quantile.IQR_{quantile:quantile.Quantile_{window:3,interpolation:%d}}",
		Midpoint,
	)
	iqr, err := NewIQR(3, AVL)
	require.NoError(t, err)

	assert.Equal(t, expectedString, iqr.String())
}

func TestIQRPush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		iqr, err := NewIQR(3, AVL)
		require.NoError(t, err)
		for i := 0.; i < 5; i++ {
			err := iqr.Push(i)
			require.NoError(t, err)
		}
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		iqr, err := NewIQR(3, AVL)
		require.NoError(t, err)

		for i := 0.; i < 3; i++ {
			err = iqr.Push(i)
			require.NoError(t, err)
		}

		// dispose the queue to simulate an error when we try to retrieve from the queue
		iqr.quantile.queue.Dispose()
		err = iqr.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		iqr, err := NewIQR(3, AVL)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to insert into the queue
		iqr.quantile.queue.Dispose()
		val := 3.
		err = iqr.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestIQRValue(t *testing.T) {
	t.Run("pass: returns IQR", func(t *testing.T) {
		iqr, err := NewIQR(3, AVL)
		require.NoError(t, err)
		for i := 0.; i < 6; i++ {
			err := iqr.Push(i)
			require.NoError(t, err)
		}

		value, err := iqr.Value()
		require.NoError(t, err)

		assert.Equal(t, 1., value)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		iqr, err := NewIQR(3, AVL)
		require.NoError(t, err)

		_, err = iqr.Value()
		testutil.ContainsError(t, err, "no values seen yet")
	})
}

func TestIQRClear(t *testing.T) {
	iqr, err := NewIQR(3, AVL)
	require.NoError(t, err)

	for i := 0.; i < 10; i++ {
		err = iqr.Push(i * i)
		require.NoError(t, err)
	}

	iqr.Clear()
	assert.Equal(t, uint64(0), iqr.quantile.queue.Len())
	assert.Equal(t, 0, iqr.quantile.statistic.Size())
}
