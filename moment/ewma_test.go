package moment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMA(t *testing.T) {
	ewma := NewEWMA(0.3)
	testutil.Approx(t, 0.3, ewma.decay)
}

type EWMAPushSuite struct {
	suite.Suite
	ewma *EWMA
}

func TestEWMAPushSuite(t *testing.T) {
	suite.Run(t, &EWMAPushSuite{})
}

func (s *EWMAPushSuite) SetupTest() {
	s.ewma = NewEWMA(0.3)
	err := Init(s.ewma)
	s.Require().NoError(err)
}

func (s *EWMAPushSuite) TestPushSuccess() {
	err := s.ewma.Push(3.)
	s.NoError(err)
}

func (s *EWMAPushSuite) TestPushFailOnNullCore() {
	ewma := NewEWMA(0.3)
	err := ewma.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

type EWMAValueSuite struct {
	suite.Suite
	ewma *EWMA
}

func TestEWMAValueSuite(t *testing.T) {
	suite.Run(t, &EWMAValueSuite{})
}

func (s *EWMAValueSuite) SetupTest() {
	s.ewma = NewEWMA(0.3)
	err := Init(s.ewma)
	s.Require().NoError(err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := s.ewma.Push(x)
		s.Require().NoError(err)
	}
}

func (s *EWMAValueSuite) TestValueSuccess() {
	value, err := s.ewma.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 4.71, value)
}

func (s *EWMAValueSuite) TestValueFailOnNullCore() {
	ewma := NewEWMA(0.3)
	_, err := ewma.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMAValueSuite) TestValueFailIfNoValuesSeen() {
	ewma := NewEWMA(0.3)
	err := Init(ewma)
	s.Require().NoError(err)

	_, err = ewma.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMAClear(t *testing.T) {
	ewma := NewEWMA(0.3)
	err := Init(ewma)
	require.NoError(t, err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := ewma.Push(x)
		require.NoError(t, err)
	}

	ewma.Clear()
	assert.Equal(t, int(0), ewma.core.count)
}

func TestEWMAString(t *testing.T) {
	ewma := NewEWMA(0.3)
	expectedString := "moment.EWMA_{decay:0.3}"
	assert.Equal(t, expectedString, ewma.String())
}
