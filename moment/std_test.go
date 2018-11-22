package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

func TestStd(t *testing.T) {
	std, err := NewStd()
	require.NoError(t, err)

	stream.TestData(std)

	value, err := std.Value()
	require.NoError(t, err)

	stream.Approx(t, math.Sqrt(7.), value)
}
