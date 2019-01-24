package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewMean(t *testing.T) {
	mean := NewMean(3)
	assert.Equal(t, 3, mean.window)
}

type MeanPushSuite struct {
	suite.Suite
	mean *Mean
}

func TestMeanPushSuite(t *testing.T) {
	suite.Run(t, &MeanPushSuite{})
}

func (s *MeanPushSuite) SetupTest() {
	s.mean = NewMean(3)
	err := Init(s.mean)
	s.Require().NoError(err)
}

func (s *MeanPushSuite) TestPushSuccess() {
	err := s.mean.Push(3.)
	s.NoError(err)
}

func (s *MeanPushSuite) TestPushFailOnNullCore() {
	mean := NewMean(3)
	err := mean.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *MeanPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.mean.core.queue.Dispose()

	err := s.mean.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *MeanPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.mean.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.mean.core.queue.Dispose()

	err := s.mean.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

type MeanValueSuite struct {
	suite.Suite
	mean *Mean
}

func TestMeanValueSuite(t *testing.T) {
	suite.Run(t, &MeanValueSuite{})
}

func (s *MeanValueSuite) SetupTest() {
	s.mean = NewMean(3)
	err := Init(s.mean)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.mean.Push(x)
		s.Require().NoError(err)
	}
}

func (s *MeanValueSuite) TestValueSuccess() {
	value, err := s.mean.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 5, value)
}

func (s *MeanValueSuite) TestValueFailOnNullCore() {
	mean := NewMean(3)
	_, err := mean.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *MeanValueSuite) TestValueFailIfNoValuesSeen() {
	mean := NewMean(3)
	err := Init(mean)
	s.Require().NoError(err)

	_, err = mean.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestMeanClear(t *testing.T) {
	mean := NewMean(3)
	err := Init(mean)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := mean.Push(x)
		require.NoError(t, err)
	}

	mean.Clear()
	assert.Equal(t, int(0), mean.core.count)
	assert.Equal(t, uint64(0), mean.core.queue.Len())
}

func TestMeanValue(t *testing.T) {
	mean := NewMean(3)
	expectedString := "moment.Mean_{window:3}"
	assert.Equal(t, expectedString, mean.String())
}
