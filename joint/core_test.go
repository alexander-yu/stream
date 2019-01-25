package joint

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

type mockMetric struct {
	vals [][]float64
	core *Core
}

func (m *mockMetric) SetCore(c *Core) {
	m.core = c
}

func (m *mockMetric) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			{2, 2},
		},
		Window: stream.IntPtr(3),
	}
}

func (m *mockMetric) String() string {
	return ""
}

func (m *mockMetric) Push(xs ...float64) error {
	m.vals = append(m.vals, xs)
	err := m.core.Push(xs...)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}

	return nil
}

func (m *mockMetric) Value() (float64, error) {
	return 0, nil
}

func (m *mockMetric) Clear() {
	m.core.Clear()
	m.vals = nil
}

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
		err := s.metric.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CorePushSuite) TestPushSuccess() {
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

	s.Equal(len(expectedSums), len(s.metric.core.sums))
	for hash, expectedSum := range expectedSums {
		actualSum := s.metric.core.sums[hash]
		testutil.Approx(s.T(), expectedSum, actualSum)
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
			testutil.Approx(s.T(), expectedVals[i][j], s.metric.vals[i][j])
		}
	}
}

func (s *CorePushSuite) TestPushSuccessForWindow1() {
	// this time, set a window of 1; the Core should really just keep the
	// most recent value. This is to test the case where we should clear out
	// any stats upon removing the last item from the queue, which only happens
	// in the special case of the queue having a size of 1.
	core, err := NewCore(&CoreConfig{
		Sums:   SumsConfig{{2, 2}},
		Window: stream.IntPtr(1),
	})
	s.Require().NoError(err)

	err = core.Push(1., 1.)
	s.Require().NoError(err)

	err = core.Push(2., 2.)
	s.Require().NoError(err)

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

	s.Equal(len(expectedSums), len(core.sums))
	for hash, expectedSum := range expectedSums {
		actualSum := core.sums[hash]
		testutil.Approx(s.T(), expectedSum, actualSum)
	}
}

func (s *CorePushSuite) TestPushFailOnQueueInsertionFailure() {
	metric := &mockMetric{}
	err := Init(metric)
	s.Require().NoError(err)

	// dispose the queue to simulate an error when we try to retrieve from the queue
	metric.core.queue.Dispose()
	vals := []float64{3., 3.}
	err = metric.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %v to queue", vals))
}

func (s *CorePushSuite) TestPushFailOnQueueRetrievalFailure() {
	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.metric.core.queue.Dispose()
	err := s.metric.Push(3., 3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

func (s *CorePushSuite) TestPushFailOnWrongNumberOfValues() {
	metric := &mockMetric{}
	err := Init(metric)
	s.Require().NoError(err)

	vals := []float64{3., 4., 5.}
	err = metric.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"tried to push %d values when core is tracking %d variables",
		len(vals),
		len(metric.core.means),
	))
}

func TestClear(t *testing.T) {
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x, x*x)
		require.NoError(t, err)
	}

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
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x, x*x)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, m.core.Count())
}

type CoreMeanSuite struct {
	suite.Suite
	metric *mockMetric
}

func TestCoreMeanSuite(t *testing.T) {
	suite.Run(t, &CoreMeanSuite{})
}

func (s *CoreMeanSuite) SetupTest() {
	s.metric = &mockMetric{}
	err := Init(s.metric)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.metric.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CoreMeanSuite) TestMeanSuccess() {
	mean, err := s.metric.core.Mean(0)
	s.Require().NoError(err)

	testutil.Approx(s.T(), 5., mean)
}

func (s *CoreMeanSuite) TestMeanFailIfNoValuesSeen() {
	core, err := NewCore(&CoreConfig{
		Vars: stream.IntPtr(2),
	})
	s.Require().NoError(err)

	_, err = core.Mean(0)
	s.EqualError(err, "no values seen yet")
}

func (s *CoreMeanSuite) TestMeanFailForInvalidVariable() {
	_, err := s.metric.core.Mean(-1)
	s.EqualError(err, fmt.Sprintf("%d is not a tracked variable", -1))

	_, err = s.metric.core.Mean(2)
	s.EqualError(err, fmt.Sprintf("%d is not a tracked variable", 2))
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
		err := s.metric.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CoreSumSuite) TestSumSuccess() {
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
		sum, err := s.metric.core.Sum(tuple...)
		s.Require().NoError(err)
		testutil.Approx(s.T(), expectedSums[tuple.hash()], sum)
	})
}

func (s *CoreSumSuite) TestFailIfNoValuesSeen() {
	core, err := NewCore(&CoreConfig{
		Vars: stream.IntPtr(2),
	})
	s.Require().NoError(err)

	_, err = core.Sum(2, 0)
	s.EqualError(err, "no values seen yet")
}

func (s *CoreSumSuite) TestFailForUntrackedSum() {
	tuple := Tuple{3, 2}
	_, err := s.metric.core.Sum(tuple...)
	s.EqualError(err, fmt.Sprintf("%v is not a tracked power sum", tuple))
}

func TestLock(t *testing.T) {
	m := &mockMetric{}
	err := Init(m)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := m.Push(x, x*x)
		require.NoError(t, err)
	}

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
