package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/quantile/skiplist"
	testutil "github.com/alexander-yu/stream/util/test"
)

func TestImplOption(t *testing.T) {
	t.Run("fail: invalid order.Option is invalid", func(t *testing.T) {
		err := ImplOption(SkipList, skiplist.ProbabilityOption(-1))(&Quantile{})
		testutil.ContainsError(t, err, "error setting Impl")
	})

	t.Run("fail: unsupported Impl is invalid", func(t *testing.T) {
		impl := Impl(-1)
		err := ImplOption(impl)(&Quantile{})
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set invalid Impl %d", impl))
	})

	t.Run("pass: valid ImplOption is valid", func(t *testing.T) {
		err := ImplOption(AVL)(&Quantile{})
		require.NoError(t, err)
	})
}

func TestInterpolationOption(t *testing.T) {
	t.Run("fail: unsupported Interpolation is invalid", func(t *testing.T) {
		i := Interpolation(-1)
		err := InterpolationOption(i)(&Quantile{})
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set invalid Interpolation %d", i))
	})

	t.Run("pass: valid InterpolationOption is valid", func(t *testing.T) {
		err := InterpolationOption(Linear)(&Quantile{})
		require.NoError(t, err)
	})
}
