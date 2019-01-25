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

func TestNewCorrelation(t *testing.T) {
	correlation := NewCorrelation(3)
	assert.Equal(t, 3, correlation.window)
}

type CorrelationPushSuite struct {
	suite.Suite
	correlation *Correlation
}

func TestCorrelationPushSuite(t *testing.T) {
	suite.Run(t, &CorrelationPushSuite{})
}

func (s *CorrelationPushSuite) SetupTest() {
	s.correlation = NewCorrelation(3)
	err := Init(s.correlation)
	s.Require().NoError(err)
}

func (s *CorrelationPushSuite) TestPushSuccess() {
	err := s.correlation.Push(3., 9.)
	s.NoError(err)
}

func (s *CorrelationPushSuite) TestPushFailOnNullCore() {
	correlation := NewCorrelation(3)
	err := correlation.Push(0., 0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CorrelationPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.correlation.core.queue.Dispose()

	err := s.correlation.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CorrelationPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.correlation.Push(x, x*x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.correlation.core.queue.Dispose()

	err := s.correlation.Push(3., 9.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *CorrelationPushSuite) TestPushFailOnWrongNumberOfValues() {
	correlation := NewCorrelation(3)
	err := Init(correlation)
	s.Require().NoError(err)

	vals := []float64{3.}
	err = correlation.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Correlation expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))

	vals = []float64{3., 9., 27.}
	err = correlation.Push(vals...)
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Correlation expected 2 arguments: got %d (%v)",
		len(vals),
		vals,
	))
}

type CorrelationValueSuite struct {
	suite.Suite
	correlation *Correlation
}

func TestCorrelationValueSuite(t *testing.T) {
	suite.Run(t, &CorrelationValueSuite{})
}

func (s *CorrelationValueSuite) SetupTest() {
	s.correlation = NewCorrelation(3)
	err := Init(s.correlation)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.correlation.Push(x, x*x)
		s.Require().NoError(err)
	}
}

func (s *CorrelationValueSuite) TestValueSuccess() {
	value, err := s.correlation.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 158./math.Sqrt(14.*5378./3.), value)
}

func (s *CorrelationValueSuite) TestValueFailOnNullCore() {
	correlation := NewCorrelation(3)
	_, err := correlation.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *CorrelationValueSuite) TestValueFailIfNoValuesSeen() {
	correlation := NewCorrelation(3)
	err := Init(correlation)
	s.Require().NoError(err)

	_, err = correlation.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestCorrelationClear(t *testing.T) {
	correlation := NewCorrelation(3)
	err := Init(correlation)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := correlation.Push(x, x*x)
		require.NoError(t, err)
	}

	correlation.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		2:  0.,
		31: 0.,
		32: 0.,
		62: 0.,
	}
	assert.Equal(t, expectedSums, correlation.core.sums)
	assert.Equal(t, expectedSums, correlation.core.newSums)
	assert.Equal(t, 0, correlation.core.count)
	assert.Equal(t, uint64(0), correlation.core.queue.Len())
}

func TestCorrelationString(t *testing.T) {
	correlation := NewCorrelation(3)
	expectedString := "joint.Correlation_{window:3}"
	assert.Equal(t, expectedString, correlation.String())
}
