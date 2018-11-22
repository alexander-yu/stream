package moment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

func TestMoment(t *testing.T) {
	moment, err := NewMoment(2)
	require.NoError(t, err)

	stream.TestData(moment)

	value, err := moment.Value()
	require.NoError(t, err)

	stream.Approx(t, 7, value)
}
