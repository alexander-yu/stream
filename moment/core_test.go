package moment

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

type mockMetric struct {
	vals []float64
	core *Core
}

func (m *mockMetric) SetCore(c *Core) {
	m.core = c
}

func (m *mockMetric) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			1: true,
			2: true,
			3: true,
			4: true,
		},
		Window: stream.IntPtr(3),
	}
}

func (m *mockMetric) String() string {
	return "mockMetric"
}

func (m *mockMetric) Push(x float64) error {
	m.vals = append(m.vals, x)
	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}

	return nil
}

func (m *mockMetric) Value() (float64, error) {
	return 0, nil
}

func (m *mockMetric) Clear() {
	m.vals = nil
	m.core.Clear()
}

type invalidMetric struct {
	coreSet bool
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

func (m *invalidMetric) SetCore(c *Core) {
	m.coreSet = true
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

		assert.Equal(t, 3, core.window)
		assert.Equal(t, make([]float64, 6), core.sums)
		assert.Equal(t, uint64(0), core.queue.Len())
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

type CorePushSuite struct {
	suite.Suite
	metric *mockMetric
}

func TestCorePushSuite(t *testing.T) {
	suite.Run(t, &CorePushSuite{})
}

func (s *CorePushSuite) SetupTest() {
	s.metric = &mockMetric{}
	err := Init(s.metric)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.metric.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CorePushSuite) TestPushSuccess() {
	expectedSums := []float64{0., 0., 14., 18., 98.}

	s.Equal(len(expectedSums), len(s.metric.core.sums))
	for k, expectedSum := range expectedSums {
		actualSum := s.metric.core.sums[k]
		testutil.Approx(s.T(), expectedSum, actualSum)
	}

	// Check that Push also pushes the value to the metric
	expectedVals := []float64{1., 2., 3., 4., 8.}
	for i := range expectedVals {
		testutil.Approx(s.T(), expectedVals[i], s.metric.vals[i])
	}
}

func (s *CorePushSuite) TestPushSuccessForWindow1() {
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
	s.Require().NoError(err)

	err = core.Push(1.)
	s.Require().NoError(err)

	err = core.Push(2.)
	s.Require().NoError(err)

	expectedSums := []float64{0., 0., 0., 0., 0.}
	s.Equal(len(expectedSums), len(core.sums))
	for k, expectedSum := range expectedSums {
		actualSum := core.sums[k]
		testutil.Approx(s.T(), expectedSum, actualSum)
	}
}

func (s *CorePushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.metric.core.queue.Dispose()
	err := s.metric.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

func (s *CorePushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.metric.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.metric.core.queue.Dispose()
	err := s.metric.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

func TestClear(t *testing.T) {
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x)
		require.NoError(t, err)
	}

	m.core.Clear()

	expectedSums := []float64{0, 0, 0, 0, 0}
	assert.Equal(t, expectedSums, m.core.sums)
	assert.Equal(t, float64(0), m.core.mean)
	assert.Equal(t, int(0), m.core.count)
	assert.Equal(t, uint64(0), m.core.queue.Len())
}

func TestCount(t *testing.T) {
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, m.core.Count())
}

type CoreMeanSuite struct {
	suite.Suite
	metric *mockMetric
}

func TestCoreMeanSuite(t *testing.T) {
	suite.Run(t, &CorePushSuite{})
}

func (s *CoreMeanSuite) SetupTest() {
	s.metric = &mockMetric{}
	err := Init(s.metric)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.metric.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CoreMeanSuite) TestMeanSuccess() {
	mean, err := s.metric.core.Mean()
	s.Require().NoError(err)

	testutil.Approx(s.T(), 5., mean)
}

func (s *CoreMeanSuite) TestMeanFailIfNoValuesSeen() {
	core, err := NewCore(&CoreConfig{})
	s.Require().NoError(err)

	_, err = core.Mean()
	s.EqualError(err, "no values seen yet")
}

type CoreSumSuite struct {
	suite.Suite
	metric *mockMetric
}

func TestCoreSumSuite(t *testing.T) {
	suite.Run(t, &CoreSumSuite{})
}

func (s *CoreSumSuite) SetupTest() {
	s.metric = &mockMetric{}
	err := Init(s.metric)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.metric.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CoreSumSuite) TestSumSuccess() {
	expectedSums := []float64{0., 0., 14., 18., 98.}

	for i := 1; i <= 4; i++ {
		sum, err := s.metric.core.Sum(i)
		s.Require().NoError(err)
		testutil.Approx(s.T(), expectedSums[i], sum)
	}
}

func (s *CoreSumSuite) TestFailIfNoValuesSeen() {
	core, err := NewCore(&CoreConfig{})
	s.Require().NoError(err)

	_, err = core.Sum(1)
	s.EqualError(err, "no values seen yet")
}

func (s *CoreSumSuite) TestFailForUntrackedSum() {
	_, err := s.metric.core.Sum(10)
	s.EqualError(err, "10 is not a tracked power sum")
}

func TestLock(t *testing.T) {
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x)
		require.NoError(t, err)
	}

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
