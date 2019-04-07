package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewStd(t *testing.T) {
	std := NewStd(3)
	assert.Equal(t, New(2, 3), std.variance)
}

func TestNewGlobalStd(t *testing.T) {
	std := NewStd(0)
	globalStd := NewGlobalStd()
	assert.Equal(t, std, globalStd)
}

type StdPushSuite struct {
	suite.Suite
	std *Std
}

func TestStdPushSuite(t *testing.T) {
	suite.Run(t, &StdPushSuite{})
}

func (s *StdPushSuite) SetupTest() {
	s.std = NewStd(3)
	err := Init(s.std)
	s.Require().NoError(err)
}

func (s *StdPushSuite) TestPushSuccess() {
	err := s.std.Push(3.)
	s.NoError(err)
}

func (s *StdPushSuite) TestPushFailOnNullCore() {
	std := NewStd(3)
	err := std.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *StdPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.std.variance.core.queue.Dispose()

	err := s.std.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *StdPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.std.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.std.variance.core.queue.Dispose()

	err := s.std.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

type StdValueSuite struct {
	suite.Suite
	std *Std
}

func TestStdValueSuite(t *testing.T) {
	suite.Run(t, &StdValueSuite{})
}

func (s *StdValueSuite) SetupTest() {
	s.std = NewStd(3)
	err := Init(s.std)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.std.Push(x)
		s.Require().NoError(err)
	}
}

func (s *StdValueSuite) TestValueSuccess() {
	value, err := s.std.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), math.Sqrt(7.), value)
}

func (s *StdValueSuite) TestValueFailOnNullCore() {
	std := NewStd(3)
	_, err := std.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *StdValueSuite) TestValueFailIfNoValuesSeen() {
	std := NewStd(3)
	err := Init(std)
	s.Require().NoError(err)

	_, err = std.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestStdClear(t *testing.T) {
	std := NewStd(3)
	err := Init(std)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := std.Push(x)
		require.NoError(t, err)
	}

	std.Clear()
	expectedSums := []float64{0, 0, 0}
	assert.Equal(t, expectedSums, std.variance.core.sums)
	assert.Equal(t, int(0), std.variance.core.count)
	assert.Equal(t, uint64(0), std.variance.core.queue.Len())
}

func TestStdString(t *testing.T) {
	std := NewStd(3)
	expectedString := "moment.Std_{window:3}"
	assert.Equal(t, expectedString, std.String())
}
