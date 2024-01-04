package minmax

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewMin(t *testing.T) {
	t.Run("pass: valid Min is valid", func(t *testing.T) {
		min, err := NewMin(3)
		require.NoError(t, err)
		assert.Equal(t, 3, min.window)
		assert.Equal(t, uint64(0), min.queue.Len())
		assert.Equal(t, 0, min.deque.Len())
		assert.Equal(t, math.Inf(1), min.min)
		assert.Equal(t, 0, min.count)
	})

	t.Run("fail: negative window returns error", func(t *testing.T) {
		_, err := NewMin(-1)
		testutil.ContainsError(t, err, "-1 is a negative window")
	})
}

func TestNewGlobalMin(t *testing.T) {
	min, err := NewMin(0)
	require.NoError(t, err)

	globalMin := NewGlobalMin()
	require.NoError(t, err)

	assert.Equal(t, min, globalMin)
}

func TestMinString(t *testing.T) {
	expectedString := "minmax.Min_{window:3}"
	min, err := NewMin(3)
	require.NoError(t, err)
	assert.Equal(t, expectedString, min.String())
}

type MinPushSuite struct {
	suite.Suite
	windowMin *Min
	globalMin *Min
}

func TestMinPushSuite(t *testing.T) {
	suite.Run(t, &MinPushSuite{})
}

func (s *MinPushSuite) SetupTest() {
	var err error
	s.globalMin, err = NewMin(0)
	s.Require().NoError(err)
	s.windowMin, err = NewMin(5)
	s.Require().NoError(err)
}

func (s *MinPushSuite) TestPushGlobalSuccess() {
	for i := 3; i > 0; i-- {
		err := s.globalMin.Push(float64(i))
		s.Require().NoError(err)

		s.Equal(4-i, s.globalMin.count)
		testutil.Approx(s.T(), float64(i), s.globalMin.min)
	}
}

func (s *MinPushSuite) TestPushWindowSuccess() {
	vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
	maxes := []float64{9, 4, 4, 1, 1, 1, 1, 1, 2, 2}

	for i, val := range vals {
		err := s.windowMin.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], s.windowMin.deque.Front())
	}

	// reset test
	s.SetupTest()
	vals = []float64{1, 2, 3, 4, 5, 6, 7}
	maxes = []float64{1, 1, 1, 1, 1, 2, 3}

	for i, val := range vals {
		err := s.windowMin.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], s.windowMin.deque.Front())
	}

	// reset test
	s.SetupTest()
	vals = []float64{7, 6, 5, 4, 3, 2, 1}
	maxes = []float64{7, 6, 5, 4, 3, 2, 1}

	for i, val := range vals {
		err := s.windowMin.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], s.windowMin.deque.Front())
	}
}

func (s *MinPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.windowMin.queue.Dispose()
	val := 3.
	err := s.windowMin.Push(val)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %f to queue", val))
}

func (s *MinPushSuite) TestPushFailOnQueueRetrievalFailure() {
	for i := 0.; i < 5; i++ {
		err := s.windowMin.Push(i)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.windowMin.queue.Dispose()
	err := s.windowMin.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

type MinValueSuite struct {
	suite.Suite
	windowMin *Min
	globalMin *Min
}

func TestMinValueSuite(t *testing.T) {
	suite.Run(t, &MinValueSuite{})
}

func (s *MinValueSuite) SetupTest() {
	var err error
	s.globalMin, err = NewMin(0)
	s.Require().NoError(err)
	s.windowMin, err = NewMin(5)
	s.Require().NoError(err)

	vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
	for _, val := range vals {
		err = s.globalMin.Push(val)
		s.Require().NoError(err)
		err = s.windowMin.Push(val)
		s.Require().NoError(err)
	}
}

func (s *MinValueSuite) TestValueGlobalSuccess() {
	val, err := s.globalMin.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 1., val)
}

func (s *MinValueSuite) TestValueWindowSuccess() {
	val, err := s.windowMin.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 2., val)
}

func (s *MinValueSuite) TestValueFailIfNoValuesSeen() {
	max, err := NewMin(3)
	s.Require().NoError(err)

	_, err = max.Value()
	assert.EqualError(s.T(), err, "no values seen yet")
}

func TestMinClear(t *testing.T) {
	min, err := NewMin(3)
	require.NoError(t, err)

	for i := 0.; i < 3; i++ {
		err = min.Push(i)
		require.NoError(t, err)
	}

	min.Clear()
	assert.Equal(t, 0, min.count)
	assert.Equal(t, math.Inf(1), min.min)
	assert.Equal(t, uint64(0), min.queue.Len())
	assert.Equal(t, 0, min.deque.Len())
}
