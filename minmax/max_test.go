package minmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestNewMax(t *testing.T) {
	t.Run("pass: returns a Max", func(t *testing.T) {
		_, err := NewMax(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewMax(-1)
		assert.EqualError(t, err, "-1 is a negative window")
	})
}

func TestMax(t *testing.T) {
	t.Run("pass: returns running global maximum for a window of 0", func(t *testing.T) {
		max, err := NewMax(0)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 9., val)
	})

	t.Run("pass: returns running maximum for a provided window", func(t *testing.T) {
		max, err := NewMax(5)
		require.NoError(t, err)

		vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
		for _, val := range vals {
			err = max.Push(val)
			require.NoError(t, err)
		}

		val, err := max.Value()
		require.NoError(t, err)
		testutil.Approx(t, 5., val)
	})

	t.Run("fail: if no values seen, return error", func(t *testing.T) {
		max, err := NewMax(3)
		require.NoError(t, err)

		_, err = max.Value()
		assert.EqualError(t, err, "no values seen yet")
	})
}
