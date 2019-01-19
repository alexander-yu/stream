package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

type invalidMetric struct {
	subscribed bool
}

func (m *invalidMetric) String() string {
	return ""
}

func (m *invalidMetric) Push(x float64) error {
	return nil
}

func (m *invalidMetric) Value() (float64, error) {
	return 0, nil
}

func (m *invalidMetric) Clear() {}

func (m *invalidMetric) Subscribe(c *Core) {
	m.subscribed = true
}

func (m *invalidMetric) Config() *CoreConfig {
	return &CoreConfig{Sums: SumsConfig{-1: true}}
}

func TestNewCore(t *testing.T) {
	t.Run("fail: invalid config returns error", func(t *testing.T) {
		_, err := NewCore(&CoreConfig{
			Sums: SumsConfig{
				-1: true,
			},
			Window: stream.IntPtr(3),
		})
		testutil.ContainsError(t, err, "error validating config")
	})

	t.Run("pass: valid config returns Core", func(t *testing.T) {
		core, err := NewCore(&CoreConfig{
			Sums: SumsConfig{
				1: true,
				5: true,
			},
			Window: stream.IntPtr(3),
		})
		require.NoError(t, err)

		assert.Equal(t, uint64(3), core.window)
		assert.Equal(t, make([]float64, 6), core.sums)
		assert.Equal(t, uint64(0), core.queue.Len())
	})
}

func TestInit(t *testing.T) {
	t.Run("fail: invalid config returns error", func(t *testing.T) {
		metric := &invalidMetric{}
		err := Init(metric)

		testutil.ContainsError(t, err, "error creating Core")
		assert.False(t, metric.subscribed)
	})
}

func TestPush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		expectedSums := []float64{0., 0., 14., 18., 98.}

		assert.Equal(t, len(expectedSums), len(m.core.sums))
		for k, expectedSum := range expectedSums {
			actualSum := m.core.sums[k]
			testutil.Approx(t, expectedSum, actualSum)
		}

		// Check that Push also pushes the value to the metric
		expectedVals := []float64{1., 2., 3., 4., 8.}
		for i := range expectedVals {
			testutil.Approx(t, expectedVals[i], m.vals[i])
		}
	})

	t.Run("pass: successfully pushes values for window of 1", func(t *testing.T) {
		// this time, set a window of 1; the Core should really just keep the
		// most recent value. This is to test the case where we should clear out
		// any stats upon removing the last item from the queue, which only happens
		// in the special case of the queue having a size of 1.
		core, err := NewCore(&CoreConfig{
			Sums: SumsConfig{
				1: true,
				2: true,
				3: true,
				4: true,
			},
			Window: stream.IntPtr(1),
		})
		require.NoError(t, err)

		err = core.Push(1.)
		require.NoError(t, err)

		err = core.Push(2.)
		require.NoError(t, err)

		expectedSums := []float64{0., 0., 0., 0., 0.}
		assert.Equal(t, len(expectedSums), len(core.sums))
		for k, expectedSum := range expectedSums {
			actualSum := core.sums[k]
			testutil.Approx(t, expectedSum, actualSum)
		}
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		m.core.queue.Dispose()
		err = m.Push(3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		m := newMockMetric()

		// dispose the queue to simulate an error when we try to insert into the queue
		m.core.queue.Dispose()
		val := 3.
		err := m.Push(val)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %f to queue", val))
	})
}

func TestClear(t *testing.T) {
	m := newMockMetric()
	err := testData(m)
	require.NoError(t, err)

	m.core.Clear()

	expectedSums := []float64{0, 0, 0, 0, 0}
	assert.Equal(t, expectedSums, m.core.sums)
	assert.Equal(t, float64(0), m.core.mean)
	assert.Equal(t, int(0), m.core.count)
	assert.Equal(t, uint64(0), m.core.queue.Len())
}

func TestCount(t *testing.T) {
	m := newMockMetric()
	err := testData(m)
	require.NoError(t, err)

	assert.Equal(t, 3, m.core.Count())
}

func TestMean(t *testing.T) {
	t.Run("pass: Mean returns the correct mean", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		mean, err := m.core.Mean()
		require.NoError(t, err)

		testutil.Approx(t, 5., mean)
	})

	t.Run("fail: Mean fails if no elements consumed yet", func(t *testing.T) {
		core, err := NewCore(&CoreConfig{})
		require.NoError(t, err)

		_, err = core.Mean()
		assert.EqualError(t, err, "no values seen yet")
	})
}

func TestSum(t *testing.T) {
	t.Run("pass: Sum returns the correct sum", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		expectedSums := []float64{0., 0., 14., 18., 98.}

		for i := 1; i <= 4; i++ {
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
		err := testData(m)
		require.NoError(t, err)

		_, err = m.core.Sum(10)
		assert.EqualError(t, err, "10 is not a tracked power sum")
	})
}

func TestLock(t *testing.T) {
	m := newMockMetric()
	err := testData(m)
	require.NoError(t, err)

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
	sum, err := m.core.Sum(2)
	require.NoError(t, err)
	testutil.Approx(t, 14., sum)

	// Undo RLock call
	m.core.RUnlock()

	// New Push call should now be unblocked
	<-done
	sum, err = m.core.Sum(2)
	require.NoError(t, err)
	testutil.Approx(t, 26./3., sum)
}
