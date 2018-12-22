package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/testutil"
)

func TestPush(t *testing.T) {
	m := newMockMetric()
	testData(m)

	expectedSums := map[int]float64{
		-1: 17. / 24.,
		0:  3.,
		1:  15.,
		2:  89.,
		3:  603.,
		4:  4433.,
	}

	assert.Equal(t, len(expectedSums), len(m.core.sums))
	for k, expectedSum := range expectedSums {
		actualSum, ok := m.core.sums[k]
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
	m := newMockMetric()
	testData(m)
	m.core.Clear()

	expectedSums := map[int]float64{
		-1: 0,
		0:  0,
		1:  0,
		2:  0,
		3:  0,
		4:  0,
	}
	assert.Equal(t, expectedSums, m.core.sums)
}

func TestCount(t *testing.T) {
	m := newMockMetric()
	testData(m)
	assert.Equal(t, 3, m.core.Count())
}

func TestSum(t *testing.T) {
	t.Run("pass: Sum returns the correct sum", func(t *testing.T) {
		m := newMockMetric()
		testData(m)
		expectedSums := map[int]float64{
			-1: 17. / 24.,
			0:  3.,
			1:  15.,
			2:  89.,
			3:  603.,
			4:  4433.,
		}

		for i := -1; i <= 4; i++ {
			sum, err := m.core.Sum(i)
			require.Nil(t, err)
			testutil.Approx(t, expectedSums[i], sum)
		}
	})

	t.Run("fail: Sum fails if no elements consumed yet", func(t *testing.T) {
		core, err := NewCore(&CoreConfig{})
		require.NoError(t, err)

		_, err = core.Sum(1)
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: Sum fails for untracked power sum", func(t *testing.T) {
		m := newMockMetric()
		testData(m)

		_, err := m.core.Sum(10)
		assert.EqualError(t, err, "10 is not a tracked power sum")
	})
}

func TestLock(t *testing.T) {
	m := newMockMetric()
	testData(m)
	done := make(chan bool)

	// Lock for reading
	m.core.RLock()

	// Spawn a goroutine to write; should be blocked
	go func() {
		m.core.Lock()
		defer m.core.Unlock()
		err := m.core.UnsafePush(5)
		require.NoError(t, err)
		done <- true
	}()

	// Read the sum; note that Sum also uses the RLock/RUnlock of the lock internally,
	// and this should not be blocked by the earlier RLock call
	sum, err := m.core.Sum(1)
	require.NoError(t, err)
	testutil.Approx(t, 15, sum)

	// Undo RLock call
	m.core.RUnlock()

	// New Push call should now be unblocked
	<-done
	sum, err = m.core.Sum(1)
	require.NoError(t, err)
	testutil.Approx(t, 17, sum)
}
