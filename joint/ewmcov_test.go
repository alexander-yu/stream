package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMCov(t *testing.T) {
	cov := NewEWMCov(0.3)
	assert.Equal(t, 0.3, cov.decay)
}

type EWMCovPushSuite struct {
	suite.Suite
	cov *EWMCov
}

func TestEWMCovPushSuite(t *testing.T) {
	suite.Run(t, &EWMCovPushSuite{})
}

func (s *EWMCovPushSuite) SetupTest() {
	s.cov = NewEWMCov(0.3)
	err := Init(s.cov)
	s.Require().NoError(err)
}

func (s *EWMCovPushSuite) TestPushSuccess() {
	err := s.cov.Push(3., 9.)
	s.NoError(err)
}

func (s *EWMCovPushSuite) TestPushFailOnNullCore() {
	cov := NewEWMCov(0.3)
	err := cov.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMCovPushSuite) TestPushFailOnWrongNumberOfValues() {
	vals := []float64{3.}
	err := s.cov.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"EWMCov expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))

	vals = []float64{3., 9., 27.}
	err = s.cov.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"EWMCov expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))
}

type EWMCovValueSuite struct {
	suite.Suite
	cov *EWMCov
}

func TestEWMCovValueSuite(t *testing.T) {
	suite.Run(t, &EWMCovValueSuite{})
}

func (s *EWMCovValueSuite) SetupTest() {
	s.cov = NewEWMCov(0.3)
	err := Init(s.cov)
	s.Require().NoError(err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := s.cov.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *EWMCovValueSuite) TestValueSuccess() {
	value, err := s.cov.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 53.2413, value)
}

func (s *EWMCovValueSuite) TestValueFailOnNullCore() {
	cov := NewEWMCov(0.3)
	_, err := cov.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMCovValueSuite) TestValueFailIfNoValuesSeen() {
	cov := NewEWMCov(0.3)
	err := Init(cov)
	s.Require().NoError(err)

	_, err = cov.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMCovClear(t *testing.T) {
	cov := NewEWMCov(0.3)
	err := Init(cov)
	require.NoError(t, err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := cov.Push(x, x*x)
		require.NoError(t, err)
	}

	cov.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		31: 0.,
		32: 0.,
	}
	assert.Equal(t, expectedSums, cov.core.sums)
	assert.Equal(t, expectedSums, cov.core.newSums)
	assert.Equal(t, 0, cov.core.count)
}

func TestEWMCovString(t *testing.T) {
	cov := NewEWMCov(0.3)
	expectedString := "joint.EWMCov_{decay:0.3}"
	assert.Equal(t, expectedString, cov.String())
}
