package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewEWMStd(t *testing.T) {
	std := NewEWMStd(0.3)
	assert.Equal(t, NewEWMMoment(2, 0.3), std.variance)
}

type EWMStdPushSuite struct {
	suite.Suite
	std *EWMStd
}

func TestEWMStdPushSuite(t *testing.T) {
	suite.Run(t, &EWMStdPushSuite{})
}

func (s *EWMStdPushSuite) SetupTest() {
	s.std = NewEWMStd(0.3)
	err := Init(s.std)
	s.Require().NoError(err)
}

func (s *EWMStdPushSuite) TestPushSuccess() {
	err := s.std.Push(3.)
	s.NoError(err)
}

func (s *EWMStdPushSuite) TestPushFailOnNullCore() {
	std := NewEWMStd(0.3)
	err := std.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

type EWMStdValueSuite struct {
	suite.Suite
	std *EWMStd
}

func TestEWMStdValueSuite(t *testing.T) {
	suite.Run(t, &EWMStdValueSuite{})
}

func (s *EWMStdValueSuite) SetupTest() {
	s.std = NewEWMStd(0.3)
	err := Init(s.std)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.std.Push(x)
		s.Require().NoError(err)
	}
}

func (s *EWMStdValueSuite) TestValueSuccess() {
	value, err := s.std.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), math.Sqrt(1.8758490975), value)
}

func (s *EWMStdValueSuite) TestValueFailOnNullCore() {
	std := NewEWMStd(0.3)
	_, err := std.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *EWMStdValueSuite) TestValueFailIfNoValuesSeen() {
	std := NewEWMStd(0.3)
	err := Init(std)
	s.Require().NoError(err)

	_, err = std.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestEWMStdClear(t *testing.T) {
	std := NewEWMStd(0.3)
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

func TestEWMStdString(t *testing.T) {
	std := NewEWMStd(0.3)
	expectedString := "moment.EWMStd_{decay:0.3}"
	assert.Equal(t, expectedString, std.String())
}
