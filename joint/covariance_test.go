package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewCovariance(t *testing.T) {
	covariance := NewCovariance(3)
	assert.Equal(t, 3, covariance.window)
}

type CovariancePushSuite struct {
	suite.Suite
	covariance *Covariance
}

func TestCovariancePushSuite(t *testing.T) {
	suite.Run(t, &CovariancePushSuite{})
}

func (s *CovariancePushSuite) SetupTest() {
	s.covariance = NewCovariance(3)
	err := Init(s.covariance)
	s.Require().NoError(err)
}

func (s *CovariancePushSuite) TestPushSuccess() {
	err := s.covariance.Push(3., 9.)
	s.NoError(err)
}

func (s *CovariancePushSuite) TestPushFailOnNullCore() {
	covariance := NewCovariance(3)
	err := covariance.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CovariancePushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.covariance.core.queue.Dispose()

	err := s.covariance.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CovariancePushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.covariance.Push(x, x*x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.covariance.core.queue.Dispose()

	err := s.covariance.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CovariancePushSuite) TestPushFailOnWrongNumberOfValues() {
	covariance := NewCovariance(3)
	err := Init(covariance)
	s.Require().NoError(err)

	vals := []float64{3.}
	err = covariance.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Covariance expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))

	vals = []float64{3., 9., 27.}
	err = covariance.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Covariance expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))
}

type CovarianceValueSuite struct {
	suite.Suite
	covariance *Covariance
}

func TestCovarianceValueSuite(t *testing.T) {
	suite.Run(t, &CovarianceValueSuite{})
}

func (s *CovarianceValueSuite) SetupTest() {
	s.covariance = NewCovariance(3)
	err := Init(s.covariance)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.covariance.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CovarianceValueSuite) TestValueSuccess() {
	value, err := s.covariance.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 79., value)
}

func (s *CovarianceValueSuite) TestValueFailOnNullCore() {
	covariance := NewCovariance(3)
	_, err := covariance.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CovarianceValueSuite) TestValueFailIfNoValuesSeen() {
	covariance := NewCovariance(3)
	err := Init(covariance)
	s.Require().NoError(err)

	_, err = covariance.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestCovarianceClear(t *testing.T) {
	covariance := NewCovariance(3)
	err := Init(covariance)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := covariance.Push(x, x*x)
		require.NoError(t, err)
	}

	covariance.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		31: 0.,
		32: 0.,
	}
	assert.Equal(t, expectedSums, covariance.core.sums)
	assert.Equal(t, expectedSums, covariance.core.newSums)
	assert.Equal(t, 0, covariance.core.count)
	assert.Equal(t, uint64(0), covariance.core.queue.Len())
}

func TestCovarianceString(t *testing.T) {
	covariance := NewCovariance(3)
	expectedString := "joint.Covariance_{window:3}"
	assert.Equal(t, expectedString, covariance.String())
}
