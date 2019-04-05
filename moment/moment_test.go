package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNew(t *testing.T) {
	moment := New(2, 3)
	assert.Equal(t, 2, moment.k)
	assert.Equal(t, 3, moment.window)
}

type MomentPushSuite struct {
	suite.Suite
	moment *Moment
}

func TestMomentPushSuite(t *testing.T) {
	suite.Run(t, &MomentPushSuite{})
}

func (s *MomentPushSuite) SetupTest() {
	s.moment = New(2, 3)
	err := Init(s.moment)
	s.Require().NoError(err)
}

func (s *MomentPushSuite) TestPushSuccess() {
	err := s.moment.Push(3.)
	s.NoError(err)
}

func (s *MomentPushSuite) TestPushFailOnNullCore() {
	moment := New(2, 3)
	err := moment.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *MomentPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.moment.core.queue.Dispose()

	err := s.moment.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *MomentPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.moment.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.moment.core.queue.Dispose()

	err := s.moment.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

type MomentValueSuite struct {
	suite.Suite
	moment *Moment
}

func TestMomentValueSuite(t *testing.T) {
	suite.Run(t, &MomentValueSuite{})
}

func (s *MomentValueSuite) SetupTest() {
	s.moment = New(2, 3)
	err := Init(s.moment)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.moment.Push(x)
		s.Require().NoError(err)
	}
}

func (s *MomentValueSuite) TestValueSuccess() {
	value, err := s.moment.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 7, value)
}

func (s *MomentValueSuite) TestValueFailOnNullCore() {
	moment := New(2, 3)
	_, err := moment.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *MomentValueSuite) TestValueFailIfNoValuesSeen() {
	moment := New(2, 3)
	err := Init(moment)
	s.Require().NoError(err)

	_, err = moment.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestMomentClear(t *testing.T) {
	moment := New(2, 3)
	err := Init(moment)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := moment.Push(x)
		require.NoError(t, err)
	}

	moment.Clear()
	expectedSums := []float64{0, 0, 0}
	assert.Equal(t, expectedSums, moment.core.sums)
	assert.Equal(t, int(0), moment.core.count)
	assert.Equal(t, uint64(0), moment.core.queue.Len())
}

func TestMomentString(t *testing.T) {
	moment := New(2, 3)
	expectedString := "moment.Moment_{k:2,window:3}"
	assert.Equal(t, expectedString, moment.String())
}
