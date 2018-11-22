package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	core := TestData()

	expectedSums := map[int]float64{
		-1: 13. / 12.,
		0:  3.,
		1:  9.,
		2:  29.,
		3:  99.,
		4:  353.,
	}
	assert.Equal(t, expectedSums, core.sums)
}

func TestClear(t *testing.T) {
	core := TestData()
	core.Clear()

	expectedSums := map[int]float64{
		-1: 0,
		0:  0,
		1:  0,
		2:  0,
		3:  0,
		4:  0,
	}
	assert.Equal(t, expectedSums, core.sums)
}

func TestMin(t *testing.T) {
	core := TestData()
	approx(t, 1, core.Min())
}

func TestMax(t *testing.T) {
	core := TestData()
	approx(t, 4, core.Max())
}

func TestCount(t *testing.T) {
	core := TestData()
	assert.Equal(t, 4, core.Count())
}

func TestSum(t *testing.T) {
	t.Run("pass: Sum returns the correct sum", func(t *testing.T) {
		core := TestData()
		expectedSums := map[int]float64{
			-1: 13. / 12.,
			0:  3.,
			1:  9.,
			2:  29.,
			3:  99.,
			4:  353.,
		}

		for i := -1; i <= 4; i++ {
			sum, err := core.Sum(i)
			require.Nil(t, err)
			assert.Equal(t, expectedSums[i], sum)
		}
	})

	t.Run("fail: Sum fails if no elements consumed yet", func(t *testing.T) {
		core := NewCore(&CoreConfig{})
		_, err := core.Sum(1)
		require.Error(err)
		assert.Equal(t, "no values seen yet", err.Error())
	})

	t.Run("fail: Sum fails for untracked power sum", func(t *testing.T) {
		core := NewCore(&CoreConfig{})
		_, err := core.Sum(4)
		require.Error(err)
		assert.Equal(t, "4 is not a tracked power sum", err.Error())
	})
}
