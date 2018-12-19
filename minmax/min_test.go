package minmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestNewMin(t *testing.T) {
	t.Run("pass: returns a Min", func(t *testing.T) {
		_, err := NewMin(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMin(-1)
		assert.EqualError(t, err, "-1 is a negative window")
	})
}

func TestMin(t *testing.T) {
	t.Run("pass: returns running global minimum for a window of 0", func(t *testing.T) {
		min, err := NewMin(0)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = min.Push(val)
			require.NoError(t, err)
		}

		val, err := min.Value()
		require.NoError(t, err)
		testutil.Approx(t, 1., val)
	})

	t.Run("pass: returns running minimum for a provided window", func(t *testing.T) {
		min, err := NewMin(5)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = min.Push(val)
			require.NoError(t, err)
		}

		val, err := min.Value()
		require.NoError(t, err)
		testutil.Approx(t, 2., val)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		min, err := NewMin(3)
		require.NoError(t, err)

		_, err = min.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
