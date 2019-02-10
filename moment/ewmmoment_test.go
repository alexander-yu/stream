package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMMoment(t *testing.T) {
	moment := NewEWMMoment(2, 0.3)
	assert.Equal(t, 2, moment.k)
	testutil.Approx(t, 0.3, moment.decay)
}

type EWMMomentPushSuite struct {
	suite.Suite
	moment *EWMMoment
}

func TestEWMMomentPushSuite(t *testing.T) {
	suite.Run(t, &EWMMomentPushSuite{})
}

func (s *EWMMomentPushSuite) SetupTest() {
	s.moment = NewEWMMoment(2, 0.3)
	err := Init(s.moment)
	s.Require().NoError(err)
}

func (s *EWMMomentPushSuite) TestPushSuccess() {
	err := s.moment.Push(3.)
	s.NoError(err)
}

func (s *EWMMomentPushSuite) TestPushFailOnNullCore() {
	moment := NewEWMMoment(2, 0.3)
	err := moment.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

type EWMMomentValueSuite struct {
	suite.Suite
	moment *EWMMoment
}

func TestEWMMomentValueSuite(t *testing.T) {
	suite.Run(t, &EWMMomentValueSuite{})
}

func (s *EWMMomentValueSuite) SetupTest() {
	s.moment = NewEWMMoment(2, 0.3)
	err := Init(s.moment)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.moment.Push(x)
		s.Require().NoError(err)
	}
}

func (s *EWMMomentValueSuite) TestValueSuccess() {
	value, err := s.moment.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 1.8758490975, value)
}

func (s *EWMMomentValueSuite) TestValueFailOnNullCore() {
	moment := NewEWMMoment(2, 0.3)
	_, err := moment.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMMomentValueSuite) TestValueFailIfNoValuesSeen() {
	moment := NewEWMMoment(2, 0.3)
	err := Init(moment)
	s.Require().NoError(err)

	_, err = moment.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMMomentClear(t *testing.T) {
	moment := NewEWMMoment(2, 0.3)
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
}

func TestEWMMomentString(t *testing.T) {
	moment := NewEWMMoment(2, 0.3)
	expectedString := "moment.EWMMoment_{k:2,decay:0.3}"
	assert.Equal(t, expectedString, moment.String())
}
