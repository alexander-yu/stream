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

func TestNewCorr(t *testing.T) {
	corr := NewCorr(3)
	assert.Equal(t, 3, corr.window)
}

func TestNewGlobalCorr(t *testing.T) {
	corr := NewCorr(0)
	globalCorr := NewGlobalCorr()
	assert.Equal(t, corr, globalCorr)
}

type CorrPushSuite struct {
	suite.Suite
	corr *Corr
}

func TestCorrPushSuite(t *testing.T) {
	suite.Run(t, &CorrPushSuite{})
}

func (s *CorrPushSuite) SetupTest() {
	s.corr = NewCorr(3)
	err := Init(s.corr)
	s.Require().NoError(err)
}

func (s *CorrPushSuite) TestPushSuccess() {
	err := s.corr.Push(3., 9.)
	s.NoError(err)
}

func (s *CorrPushSuite) TestPushFailOnNullCore() {
	corr := NewCorr(3)
	err := corr.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CorrPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.corr.core.queue.Dispose()

	err := s.corr.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CorrPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.corr.Push(x, x*x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.corr.core.queue.Dispose()

	err := s.corr.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CorrPushSuite) TestPushFailOnWrongNumberOfValues() {
	corr := NewCorr(3)
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

type CorrValueSuite struct {
	suite.Suite
	corr *Corr
}

func TestCorrValueSuite(t *testing.T) {
	suite.Run(t, &CorrValueSuite{})
}

func (s *CorrValueSuite) SetupTest() {
	s.corr = NewCorr(3)
	err := Init(s.corr)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.corr.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CorrValueSuite) TestValueSuccess() {
	value, err := s.corr.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 158./math.Sqrt(14.*5378./3.), value)
}

func (s *CorrValueSuite) TestValueFailOnNullCore() {
	corr := NewCorr(3)
	_, err := corr.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CorrValueSuite) TestValueFailIfNoValuesSeen() {
	corr := NewCorr(3)
	err := Init(corr)
	s.Require().NoError(err)

	_, err = corr.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestCorrClear(t *testing.T) {
	corr := NewCorr(3)
	err := Init(corr)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
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
	assert.Equal(t, uint64(0), corr.core.queue.Len())
}

func TestCorrString(t *testing.T) {
	corr := NewCorr(3)
	expectedString := "joint.Corr_{window:3}"
	assert.Equal(t, expectedString, corr.String())
}
