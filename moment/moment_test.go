package moment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"stream"
)

func TestMoment(t *testing.T) {
	moment, err := NewMoment(2)
	require.NoError(t, err)

	core := stream.TestData(moment)
	core.Push(8)

	value, err := moment.Value()
	require.NoError(t, err)

	approx(t, 7., value)
}
