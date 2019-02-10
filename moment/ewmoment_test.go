package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMoment(t *testing.T) {
	moment := NewEWMoment(2, 0.3)
	assert.Equal(t, 2, moment.k)
	testutil.Approx(t, 0.3, moment.decay)
}

type EWMomentPushSuite struct {
	suite.Suite
	moment *EWMoment
}

func TestEWMomentPushSuite(t *testing.T) {
	suite.Run(t, &EWMomentPushSuite{})
}

func (s *EWMomentPushSuite) SetupTest() {
	s.moment = NewEWMoment(2, 0.3)
	err := Init(s.moment)
	s.Require().NoError(err)
}

func (s *EWMomentPushSuite) TestPushSuccess() {
	err := s.moment.Push(3.)
	s.NoError(err)
}

func (s *EWMomentPushSuite) TestPushFailOnNullCore() {
	moment := NewEWMoment(2, 0.3)
	err := moment.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

type EWMomentValueSuite struct {
	suite.Suite
	moment *EWMoment
}

func TestEWMomentValueSuite(t *testing.T) {
	suite.Run(t, &EWMomentValueSuite{})
}

func (s *EWMomentValueSuite) SetupTest() {
	s.moment = NewEWMoment(2, 0.3)
	err := Init(s.moment)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.moment.Push(x)
		s.Require().NoError(err)
	}
}

func (s *EWMomentValueSuite) TestValueSuccess() {
	value, err := s.moment.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 1.8758490975, value)
}

func (s *EWMomentValueSuite) TestValueFailOnNullCore() {
	moment := NewEWMoment(2, 0.3)
	_, err := moment.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMomentValueSuite) TestValueFailIfNoValuesSeen() {
	moment := NewEWMoment(2, 0.3)
	err := Init(moment)
	s.Require().NoError(err)

	_, err = moment.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMomentClear(t *testing.T) {
	moment := NewEWMoment(2, 0.3)
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

func TestEWMomentString(t *testing.T) {
	moment := NewEWMoment(2, 0.3)
	expectedString := "moment.EWMoment_{k:2,decay:0.3}"
	assert.Equal(t, expectedString, moment.String())
}
