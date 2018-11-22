package moment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"stream"
)

func TestMean(t *testing.T) {
	mean, err := NewMean()
	require.NoError(err)

	core := stream.TestData(mean)

	value, err := mean.Value()
	require.NoError(t, err)

	approx(t, 3., value)
}
