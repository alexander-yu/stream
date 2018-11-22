package moment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

func TestMean(t *testing.T) {
	mean := &Mean{}

	stream.TestData(mean)

	value, err := mean.Value()
	require.NoError(t, err)

	stream.Approx(t, 5, value)
}
