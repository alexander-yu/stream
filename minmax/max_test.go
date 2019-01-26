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

func TestNewMax(t *testing.T) {
	t.Run("pass: valid Max is valid", func(t *testing.T) {
		max, err := NewMax(3)
		require.NoError(t, err)

		assert.Equal(t, 3, max.window)
		assert.Equal(t, uint64(0), max.queue.Len())
		assert.Equal(t, 0, max.deque.Len())
		assert.Equal(t, math.Inf(-1), max.max)
		assert.Equal(t, 0, max.count)
	})

	t.Run("fail: negative window returns error", func(t *testing.T) {
		_, err := NewMax(-1)
		testutil.ContainsError(t, err, "-1 is a negative window")
	})
}

func TestMaxString(t *testing.T) {
	expectedString := "minmax.Max_{window:3}"
	max, err := NewMax(3)
	require.NoError(t, err)
	assert.Equal(t, expectedString, max.String())
}

type MaxPushSuite struct {
	suite.Suite
	windowMax *Max
	globalMax *Max
}

func TestMaxPushSuite(t *testing.T) {
	suite.Run(t, &MaxPushSuite{})
}

func (s *MaxPushSuite) SetupTest() {
	var err error
	s.globalMax, err = NewMax(0)
	s.Require().NoError(err)
	s.windowMax, err = NewMax(5)
	s.Require().NoError(err)
}

func (s *MaxPushSuite) TestPushGlobalSuccess() {
	for i := 0; i < 3; i++ {
		err := s.globalMax.Push(float64(i))
		s.Require().NoError(err)

		s.Equal(i+1, s.globalMax.count)
		testutil.Approx(s.T(), float64(i), s.globalMax.max)
	}
}

func (s *MaxPushSuite) TestPushWindowSuccess() {
	vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
	maxes := []float64{9, 9, 9, 9, 9, 8, 8, 8, 8, 5}

	for i, val := range vals {
		err := s.windowMax.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], *s.windowMax.deque.Front().(*float64))
	}

	// reset test
	s.SetupTest()
	vals = []float64{1, 2, 3, 4, 5, 6, 7}
	maxes = []float64{1, 2, 3, 4, 5, 6, 7}

	for i, val := range vals {
		err := s.windowMax.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], *s.windowMax.deque.Front().(*float64))
	}

	// reset test
	s.SetupTest()
	vals = []float64{7, 6, 5, 4, 3, 2, 1}
	maxes = []float64{7, 7, 7, 7, 7, 6, 5}

	for i, val := range vals {
		err := s.windowMax.Push(val)
		s.Require().NoError(err)
		s.Equal(maxes[i], *s.windowMax.deque.Front().(*float64))
	}
}

func (s *MaxPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.windowMax.queue.Dispose()
	val := 3.
	err := s.windowMax.Push(val)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %f to queue", val))
}

func (s *MaxPushSuite) TestPushFailOnQueueRetrievalFailure() {
	for i := 0.; i < 5; i++ {
		err := s.windowMax.Push(i)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.windowMax.queue.Dispose()
	err := s.windowMax.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from queue")
}

type MaxValueSuite struct {
	suite.Suite
	windowMax *Max
	globalMax *Max
}

func TestMaxValueSuite(t *testing.T) {
	suite.Run(t, &MaxValueSuite{})
}

func (s *MaxValueSuite) SetupTest() {
	var err error
	s.globalMax, err = NewMax(0)
	s.Require().NoError(err)
	s.windowMax, err = NewMax(5)
	s.Require().NoError(err)

	vals := []float64{9, 4, 6, 1, 8, 2, 2, 5, 5, 3}
	for _, val := range vals {
		err = s.globalMax.Push(val)
		s.Require().NoError(err)
		err = s.windowMax.Push(val)
		s.Require().NoError(err)
	}
}

func (s *MaxValueSuite) TestValueGlobalSuccess() {
	val, err := s.globalMax.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 9., val)
}

func (s *MaxValueSuite) TestValueWindowSuccess() {
	val, err := s.windowMax.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 5., val)
}

func (s *MaxValueSuite) TestValueFailIfNoValuesSeen() {
	max, err := NewMax(3)
	s.Require().NoError(err)

	_, err = max.Value()
	assert.EqualError(s.T(), err, "no values seen yet")
}

func TestMaxClear(t *testing.T) {
	max, err := NewMax(3)
	require.NoError(t, err)

	for i := 0.; i < 3; i++ {
		err = max.Push(i)
		require.NoError(t, err)
	}

	max.Clear()
	assert.Equal(t, 0, max.count)
	assert.Equal(t, math.Inf(-1), max.max)
	assert.Equal(t, uint64(0), max.queue.Len())
	assert.Equal(t, 0, max.deque.Len())
}
