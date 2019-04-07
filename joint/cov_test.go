package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewCov(t *testing.T) {
	cov := NewCov(3)
	assert.Equal(t, 3, cov.window)
}

func TestNewGlobalCov(t *testing.T) {
	cov := NewCov(0)
	globalCov := NewGlobalCov()
	assert.Equal(t, cov, globalCov)
}

type CovPushSuite struct {
	suite.Suite
	cov *Cov
}

func TestCovPushSuite(t *testing.T) {
	suite.Run(t, &CovPushSuite{})
}

func (s *CovPushSuite) SetupTest() {
	s.cov = NewCov(3)
	err := Init(s.cov)
	s.Require().NoError(err)
}

func (s *CovPushSuite) TestPushSuccess() {
	err := s.cov.Push(3., 9.)
	s.NoError(err)
}

func (s *CovPushSuite) TestPushFailOnNullCore() {
	cov := NewCov(3)
	err := cov.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CovPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.cov.core.queue.Dispose()

	err := s.cov.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CovPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.cov.Push(x, x*x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.cov.core.queue.Dispose()

	err := s.cov.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CovPushSuite) TestPushFailOnWrongNumberOfValues() {
	cov := NewCov(3)
	err := Init(cov)
	s.Require().NoError(err)

	vals := []float64{3.}
	err = cov.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Cov expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))

	vals = []float64{3., 9., 27.}
	err = cov.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Cov expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))
}

type CovValueSuite struct {
	suite.Suite
	cov *Cov
}

func TestCovValueSuite(t *testing.T) {
	suite.Run(t, &CovValueSuite{})
}

func (s *CovValueSuite) SetupTest() {
	s.cov = NewCov(3)
	err := Init(s.cov)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.cov.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CovValueSuite) TestValueSuccess() {
	value, err := s.cov.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 79., value)
}

func (s *CovValueSuite) TestValueFailOnNullCore() {
	cov := NewCov(3)
	_, err := cov.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CovValueSuite) TestValueFailIfNoValuesSeen() {
	cov := NewCov(3)
	err := Init(cov)
	s.Require().NoError(err)

	_, err = cov.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestCovClear(t *testing.T) {
	cov := NewCov(3)
	err := Init(cov)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
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
	assert.Equal(t, uint64(0), cov.core.queue.Len())
}

func TestCovString(t *testing.T) {
	cov := NewCov(3)
	expectedString := "joint.Cov_{window:3}"
	assert.Equal(t, expectedString, cov.String())
}
