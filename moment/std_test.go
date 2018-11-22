package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"stream"
)

func TestStd(t *testing.T) {
	std, err := NewStd()
	require.NoError(t, err)

	core := stream.TestData(std)
	core.Push(8)

	value, err := std.Value()
	require.NoError(t, err)

	approx(t, math.Sqrt(7.), value)
}
