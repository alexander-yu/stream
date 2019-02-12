package joint

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMCorr(t *testing.T) {
	corr := NewEWMCorr(0.3)
	assert.Equal(t, 0.3, corr.decay)
}

type EWMCorrPushSuite struct {
	suite.Suite
	corr *EWMCorr
}

func TestEWMCorrPushSuite(t *testing.T) {
	suite.Run(t, &EWMCorrPushSuite{})
}

func (s *EWMCorrPushSuite) SetupTest() {
	s.corr = NewEWMCorr(0.3)
	err := Init(s.corr)
	s.Require().NoError(err)
}

func (s *EWMCorrPushSuite) TestPushSuccess() {
	err := s.corr.Push(3., 9.)
	s.NoError(err)
}

func (s *EWMCorrPushSuite) TestPushFailOnNullCore() {
	corr := NewEWMCorr(0.3)
	err := corr.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMCorrPushSuite) TestPushFailOnWrongNumberOfValues() {
	corr := NewEWMCorr(0.3)
	err := Init(corr)
	s.Require().NoError(err)

	vals := []float64{3.}
	err = corr.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Corr expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))

	vals = []float64{3., 9., 27.}
	err = corr.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Corr expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))
}

type EWMCorrValueSuite struct {
	suite.Suite
	corr *EWMCorr
}

func TestEWMCorrValueSuite(t *testing.T) {
	suite.Run(t, &EWMCorrValueSuite{})
}

func (s *EWMCorrValueSuite) SetupTest() {
	s.corr = NewEWMCorr(0.3)
	err := Init(s.corr)
	s.Require().NoError(err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := s.corr.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *EWMCorrValueSuite) TestValueSuccess() {
	value, err := s.corr.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 53.2413/math.Sqrt(4.7859*594.8691), value)
}

func (s *EWMCorrValueSuite) TestValueFailOnNullCore() {
	corr := NewEWMCorr(0.3)
	_, err := corr.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMCorrValueSuite) TestValueFailIfNoValuesSeen() {
	corr := NewEWMCorr(0.3)
	err := Init(corr)
	s.Require().NoError(err)

	_, err = corr.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMCorrClear(t *testing.T) {
	corr := NewEWMCorr(0.3)
	err := Init(corr)
	require.NoError(t, err)

	xs := []float64{3, 4, 8}
	for _, x := range xs {
		err := corr.Push(x, x*x)
		require.NoError(t, err)
	}

	corr.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		2:  0.,
		31: 0.,
		32: 0.,
		62: 0.,
	}
	assert.Equal(t, expectedSums, corr.core.sums)
	assert.Equal(t, expectedSums, corr.core.newSums)
	assert.Equal(t, 0, corr.core.count)
}

func TestEWMCorrString(t *testing.T) {
	corr := NewEWMCorr(0.3)
	expectedString := "joint.EWMCorr_{decay:0.3}"
	assert.Equal(t, expectedString, corr.String())
}
