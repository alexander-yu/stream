package moment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alexander-yu/stream"
	testutil "github.com/alexander-yu/stream/util/test"
)

type mockWrapper struct {
	core *Core
}

func (w *mockWrapper) SetCore(c *Core) {
	w.core = c
}

func (w *mockWrapper) Config() *CoreConfig {
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

type invalidWrapper struct {
	coreSet bool
}

func (w *invalidWrapper) SetCore(c *Core) {
	w.coreSet = true
}

func (w *invalidWrapper) Config() *CoreConfig {
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
		wrapper := &invalidWrapper{}
		err := Init(wrapper)

		testutil.ContainsError(t, err, "error creating Core")
		assert.False(t, wrapper.coreSet)
	})

	t.Run("pass: valid config sets Core for wrapper", func(t *testing.T) {
		wrapper := &mockWrapper{}
		err := Init(wrapper)
		require.NoError(t, err)

		assert.NotNil(t, wrapper.core)
	})
}

type CorePushSuite struct {
	suite.Suite
	wrapper *mockWrapper
}

func TestCorePushSuite(t *testing.T) {
	suite.Run(t, &CorePushSuite{})
}

func (s *CorePushSuite) SetupTest() {
	s.wrapper = &mockWrapper{}
	err := Init(s.wrapper)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.wrapper.core.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CorePushSuite) TestPushSuccess() {
	expectedSums := []float64{0., 0., 14., 18., 98.}

	s.Equal(len(expectedSums), len(s.wrapper.core.sums))
	for k, expectedSum := range expectedSums {
		actualSum := s.wrapper.core.sums[k]
		testutil.Approx(s.T(), expectedSum, actualSum)
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
	wrapper := &mockWrapper{}
	err := Init(wrapper)
	s.Require().NoError(err)

	// dispose the queue to simulate an error when we try to retrieve from the queue
	wrapper.core.queue.Dispose()
	err = wrapper.core.Push(3.)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %f to queue", 3.))
}

func (s *CorePushSuite) TestPushFailOnQueueRetrievalFailure() {
	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.wrapper.core.queue.Dispose()
	err := s.wrapper.core.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

func TestClear(t *testing.T) {
	wrapper := &mockWrapper{}
	err := Init(wrapper)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := wrapper.core.Push(x)
		require.NoError(t, err)
	}

	wrapper.core.Clear()

	expectedSums := []float64{0, 0, 0, 0, 0}
	assert.Equal(t, expectedSums, wrapper.core.sums)
	assert.Equal(t, float64(0), wrapper.core.mean)
	assert.Equal(t, int(0), wrapper.core.count)
	assert.Equal(t, uint64(0), wrapper.core.queue.Len())
}

func TestCount(t *testing.T) {
	wrapper := &mockWrapper{}
	err := Init(wrapper)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := wrapper.core.Push(x)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, wrapper.core.Count())
}

type CoreMeanSuite struct {
	suite.Suite
	wrapper *mockWrapper
}

func TestCoreMeanSuite(t *testing.T) {
	suite.Run(t, &CoreMeanSuite{})
}

func (s *CoreMeanSuite) SetupTest() {
	s.wrapper = &mockWrapper{}
	err := Init(s.wrapper)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.wrapper.core.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CoreMeanSuite) TestMeanSuccess() {
	mean, err := s.wrapper.core.Mean()
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
	wrapper *mockWrapper
}

func TestCoreSumSuite(t *testing.T) {
	suite.Run(t, &CoreSumSuite{})
}

func (s *CoreSumSuite) SetupTest() {
	s.wrapper = &mockWrapper{}
	err := Init(s.wrapper)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.wrapper.core.Push(x)
		s.Require().NoError(err)
	}
}

func (s *CoreSumSuite) TestSumSuccess() {
	expectedSums := []float64{0., 0., 14., 18., 98.}

	for i := 1; i <= 4; i++ {
		sum, err := s.wrapper.core.Sum(i)
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
	_, err := s.wrapper.core.Sum(10)
	s.EqualError(err, "10 is not a tracked power sum")
}

func TestLock(t *testing.T) {
	wrapper := &mockWrapper{}
	err := Init(wrapper)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := wrapper.core.Push(x)
		require.NoError(t, err)
	}

	done := make(chan bool)

	// Lock for reading
	wrapper.core.RLock()

	// Spawn a goroutine to write; should be blocked
	go func() {
		wrapper.core.Lock()
		defer wrapper.core.Unlock()
		err := wrapper.core.UnsafePush(5)
		require.NoError(t, err)
		done <- true
	}()

	// Read the sum; note that Sum also uses the RLock/RUnlock of the lock internally,
	// and this should not be blocked by the earlier RLock call
	sum, err := wrapper.core.Sum(2)
	require.NoError(t, err)
	testutil.Approx(t, 14., sum)

	// Undo RLock call
	wrapper.core.RUnlock()

	// New Push call should now be unblocked
	<-done
	sum, err = wrapper.core.Sum(2)
	require.NoError(t, err)
	testutil.Approx(t, 26./3., sum)
}
