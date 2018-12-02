package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestPush(t *testing.T) {
	m := &mockMetric{}
	core := TestData(m)

	expectedSums := map[int]float64{
		-1: 17. / 24.,
		0:  3.,
		1:  15.,
		2:  89.,
		3:  603.,
		4:  4433.,
	}

	assert.Equal(t, len(expectedSums), len(core.sums))
	for k, expectedSum := range expectedSums {
		actualSum, ok := core.sums[k]
		require.True(t, ok)
		testutil.Approx(t, expectedSum, actualSum)
	}

	// Check that Push also pushes the value to the metric
	expectedVals := []float64{1., 2., 3., 4., 8.}
	for i := range expectedVals {
		testutil.Approx(t, expectedVals[i], m.vals[i])
	}
}

func TestClear(t *testing.T) {
	core := TestData(&mockMetric{})
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
	core := TestData(&mockMetric{})
	testutil.Approx(t, 1, core.Min())
}

func TestMax(t *testing.T) {
	core := TestData(&mockMetric{})
	testutil.Approx(t, 8, core.Max())
}

func TestCount(t *testing.T) {
	core := TestData(&mockMetric{})
	assert.Equal(t, 3, core.Count())
}

func TestSum(t *testing.T) {
	t.Run("pass: Sum returns the correct sum", func(t *testing.T) {
		core := TestData(&mockMetric{})
		expectedSums := map[int]float64{
			-1: 17. / 24.,
			0:  3.,
			1:  15.,
			2:  89.,
			3:  603.,
			4:  4433.,
		}

		for i := -1; i <= 4; i++ {
			sum, err := core.Sum(i)
			require.Nil(t, err)
			testutil.Approx(t, expectedSums[i], sum)
		}
	})

	t.Run("fail: Sum fails if no elements consumed yet", func(t *testing.T) {
		core := NewCore(&CoreConfig{})

		_, err := core.Sum(1)
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: Sum fails for untracked power sum", func(t *testing.T) {
		core := TestData(&mockMetric{})

		_, err := core.Sum(10)
		assert.EqualError(t, err, "10 is not a tracked power sum")
	})
}
