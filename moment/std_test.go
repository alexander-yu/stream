package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/testutil"
)

func TestStd(t *testing.T) {
	std, err := NewStd()
	require.NoError(t, err)

	stream.TestData(std)

	value, err := std.Value()
	require.NoError(t, err)

	testutil.Approx(t, math.Sqrt(7.), value)
}
