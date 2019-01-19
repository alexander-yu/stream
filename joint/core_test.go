package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

type invalidMetric struct {
	coreSet bool
}

func (m *invalidMetric) String() string {
	return ""
}

func (m *invalidMetric) Push(xs ...float64) error {
	return nil
}

func (m *invalidMetric) Value() (float64, error) {
	return 0, nil
}

func (m *invalidMetric) Clear() {}

func (m *invalidMetric) SetCore(c *Core) {
	m.coreSet = true
}

func (m *invalidMetric) Config() *CoreConfig {
	return &CoreConfig{Vars: stream.IntPtr(-1)}
}

func TestNewCore(t *testing.T) {
	t.Run("fail: invalid config returns error", func(t *testing.T) {
		_, err := NewCore(&CoreConfig{Vars: stream.IntPtr(-1)})
		testutil.ContainsError(t, err, "error validating config")
	})

	t.Run("pass: valid config returns Core", func(t *testing.T) {
		config := &CoreConfig{
			Sums: SumsConfig{
				{2, 2},
				{3, 1},
			},
			Window: stream.IntPtr(2),
		}
		core, err := NewCore(config)
		require.NoError(t, err)

		assert.Equal(t, 2, core.window)
		assert.Equal(t, config.Sums, SumsConfig(core.tuples))
		assert.Equal(t, make([]float64, 2), core.means)
		assert.Equal(t, uint64(0), core.queue.Len())

		for _, tuple := range config.Sums {
			iter(tuple, false, func(xs ...int) {
				assert.Equal(t, 0., core.sums[Tuple(xs).hash()])
			})
			iter(tuple, false, func(xs ...int) {
				assert.Equal(t, 0., core.newSums[Tuple(xs).hash()])
			})
		}
	})
}

func TestInit(t *testing.T) {
	t.Run("fail: invalid config returns error", func(t *testing.T) {
		metric := &invalidMetric{}
		err := Init(metric)

		testutil.ContainsError(t, err, "error creating Core")
		assert.False(t, metric.coreSet)
	})

	t.Run("pass: valid config sets Core for metric", func(t *testing.T) {
		metric := &mockMetric{}
		err := Init(metric)
		require.NoError(t, err)

		assert.NotNil(t, metric.core)
	})
}

func TestPush(t *testing.T) {
	t.Run("pass: successfully pushes values", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		expectedSums := map[uint64]float64{
			0:  0.,
			1:  0.,
			2:  5378. / 3.,
			31: 0.,
			32: 158.,
			33: 7486. / 3.,
			62: 14.,
			63: 638. / 3.,
			64: 112538. / 9.,
		}

		assert.Equal(t, len(expectedSums), len(m.core.sums))
		for hash, expectedSum := range expectedSums {
			actualSum := m.core.sums[hash]
			testutil.Approx(t, expectedSum, actualSum)
		}

		// Check that Push also pushes the value to the metric
		expectedVals := [][]float64{
			{1., 1.},
			{2., 4.},
			{3., 9.},
			{4., 16.},
			{8., 64.},
		}
		for i := range expectedVals {
			for j := range expectedVals[0] {
				testutil.Approx(t, expectedVals[i][j], m.vals[i][j])
			}
		}
	})

	t.Run("pass: successfully pushes values for window of 1", func(t *testing.T) {
		// this time, set a window of 1; the Core should really just keep the
		// most recent value. This is to test the case where we should clear out
		// any stats upon removing the last item from the queue, which only happens
		// in the special case of the queue having a size of 1.
		core, err := NewCore(&CoreConfig{
			Sums:   SumsConfig{{2, 2}},
			Window: stream.IntPtr(1),
		})
		require.NoError(t, err)

		err = core.Push(1., 1.)
		require.NoError(t, err)

		err = core.Push(2., 2.)
		require.NoError(t, err)

		expectedSums := map[uint64]float64{
			0:  0.,
			1:  0.,
			2:  0.,
			31: 0.,
			32: 0.,
			33: 0.,
			62: 0.,
			63: 0.,
			64: 0.,
		}

		assert.Equal(t, len(expectedSums), len(core.sums))
		for hash, expectedSum := range expectedSums {
			actualSum := core.sums[hash]
			testutil.Approx(t, expectedSum, actualSum)
		}
	})

	t.Run("fail: if queue retrieval fails, return error", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		// dispose the queue to simulate an error when we try to retrieve from the queue
		m.core.queue.Dispose()
		err = m.Push(3., 3.)
		testutil.ContainsError(t, err, "error popping item from queue")
	})

	t.Run("fail: if queue insertion fails, return error", func(t *testing.T) {
		m := newMockMetric()

		// dispose the queue to simulate an error when we try to insert into the queue
		m.core.queue.Dispose()
		vals := []float64{3., 3.}
		err := m.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf("error pushing %v to queue", vals))
	})

	t.Run("fail: if values pushed does not match Vars, return error", func(t *testing.T) {
		m := newMockMetric()

		vals := []float64{3., 4., 5.}
		err := m.Push(vals...)
		testutil.ContainsError(t, err, fmt.Sprintf(
			"tried to push %d values when core is tracking %d variables",
			len(vals),
			len(m.core.means),
		))
	})
}

func TestClear(t *testing.T) {
	m := newMockMetric()
	err := testData(m)
	require.NoError(t, err)

	m.core.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		2:  0.,
		31: 0.,
		32: 0.,
		33: 0.,
		62: 0.,
		63: 0.,
		64: 0.,
	}
	assert.Equal(t, expectedSums, m.core.sums)
	assert.Equal(t, expectedSums, m.core.newSums)

	expectedMeans := []float64{0, 0}
	assert.Equal(t, expectedMeans, m.core.means)
	assert.Equal(t, 0, m.core.count)
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

		mean, err := m.core.Mean(0)
		require.NoError(t, err)

		testutil.Approx(t, 5., mean)
	})

	t.Run("fail: Mean fails if no elements consumed yet", func(t *testing.T) {
		core, err := NewCore(&CoreConfig{
			Vars: stream.IntPtr(2),
		})
		require.NoError(t, err)

		_, err = core.Mean(0)
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: Mean fails for invalid variable", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		_, err = m.core.Mean(-1)
		assert.EqualError(t, err, fmt.Sprintf("%d is not a tracked variable", -1))

		_, err = m.core.Mean(2)
		assert.EqualError(t, err, fmt.Sprintf("%d is not a tracked variable", 2))
	})
}

func TestSum(t *testing.T) {
	t.Run("pass: Sum returns the correct sum", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		expectedSums := map[uint64]float64{
			0:  0.,
			1:  0.,
			2:  5378. / 3.,
			31: 0.,
			32: 158.,
			33: 7486. / 3.,
			62: 14.,
			63: 638. / 3.,
			64: 112538. / 9.,
		}

		iter(Tuple{2, 2}, false, func(xs ...int) {
			tuple := Tuple(xs)
			sum, err := m.core.Sum(tuple...)
			require.NoError(t, err)
			testutil.Approx(t, expectedSums[tuple.hash()], sum)
		})
	})

	t.Run("fail: Sum fails if no elements consumed yet", func(t *testing.T) {
		core, err := NewCore(&CoreConfig{
			Vars: stream.IntPtr(2),
		})
		require.NoError(t, err)

		_, err = core.Sum(2, 0)
		assert.EqualError(t, err, "no values seen yet")
	})

	t.Run("fail: Sum fails for untracked power sum", func(t *testing.T) {
		m := newMockMetric()
		err := testData(m)
		require.NoError(t, err)

		tuple := Tuple{3, 2}
		_, err = m.core.Sum(tuple...)
		assert.EqualError(t, err, fmt.Sprintf("%v is not a tracked power sum", tuple))
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
		err := m.core.UnsafePush(5., 25.)
		require.NoError(t, err)
		done <- true
	}()

	// Read the sum; note that Sum also uses the RLock/RUnlock of the lock internally,
	// and this should not be blocked by the earlier RLock call
	sum, err := m.core.Sum(2, 0)
	require.NoError(t, err)
	testutil.Approx(t, 14., sum)

	// Undo RLock call
	m.core.RUnlock()

	// New Push call should now be unblocked
	<-done
	sum, err = m.core.Sum(2, 0)
	require.NoError(t, err)
	testutil.Approx(t, 26./3., sum)
}
