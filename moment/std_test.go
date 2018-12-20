package moment

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestNewStd(t *testing.T) {
	t.Run("pass: returns an Std", func(t *testing.T) {
		_, err := NewStd(3)
		assert.NoError(t, err)
	})

	t.Run("fail: negative window is invalid", func(t *testing.T) {
		_, err := NewStd(-1)
		testutil.ContainsError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})
}

func TestStd(t *testing.T) {
	t.Run("pass: returns the standard deviation", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		testData(std)

		value, err := std.Value()
		require.NoError(t, err)

		testutil.Approx(t, math.Sqrt(7.), value)
	})

	t.Run("fail: error if no values are seen", func(t *testing.T) {
		std, err := NewStd(3)
		require.NoError(t, err)

		_, err = std.Value()
		assert.EqualError(t, err, "error retrieving 2nd moment: no values seen yet")
	})
}
