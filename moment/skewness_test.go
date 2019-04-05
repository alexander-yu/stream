package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewSkewness(t *testing.T) {
	window := 3
	skewness := NewSkewness(window)
	assert.Equal(t, &Skewness{
		variance: New(2, window),
		moment3:  New(3, window),
		config: &CoreConfig{
			Sums: SumsConfig{
				2: true,
				3: true,
			},
			Window: &window,
		},
	}, skewness)
}

type SkewnessPushSuite struct {
	suite.Suite
	skewness *Skewness
}

func TestSkewnessPushSuite(t *testing.T) {
	suite.Run(t, &SkewnessPushSuite{})
}

func (s *SkewnessPushSuite) SetupTest() {
	s.skewness = NewSkewness(3)
	err := Init(s.skewness)
	s.Require().NoError(err)
}

func (s *SkewnessPushSuite) TestPushSuccess() {
	err := s.skewness.Push(3.)
	s.NoError(err)
}

func (s *SkewnessPushSuite) TestPushFailOnNullCore() {
	skewness := NewSkewness(3)
	err := skewness.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *SkewnessPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.skewness.core.queue.Dispose()

	err := s.skewness.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *SkewnessPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.skewness.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.skewness.core.queue.Dispose()

	err := s.skewness.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

type SkewnessValueSuite struct {
	suite.Suite
	skewness *Skewness
}

func TestSkewnessValueSuite(t *testing.T) {
	suite.Run(t, &SkewnessValueSuite{})
}

func (s *SkewnessValueSuite) SetupTest() {
	s.skewness = NewSkewness(3)
	err := Init(s.skewness)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.skewness.Push(x)
		s.Require().NoError(err)
	}
}

func (s *SkewnessValueSuite) TestValueSuccess() {
	value, err := s.skewness.Value()
	s.Require().NoError(err)

	adjust := 3.
	moment := 9.
	variance := 7.

	testutil.Approx(s.T(), adjust*moment/math.Pow(variance, 1.5), value)
}

func (s *SkewnessValueSuite) TestValueFailOnNullCore() {
	skewness := NewSkewness(3)
	_, err := skewness.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *SkewnessValueSuite) TestValueFailIfNoValuesSeen() {
	skewness := NewSkewness(3)
	err := Init(skewness)
	s.Require().NoError(err)

	_, err = skewness.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestSkewnessClear(t *testing.T) {
	skewness := NewSkewness(3)
	err := Init(skewness)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := skewness.Push(x)
		require.NoError(t, err)
	}

	skewness.Clear()
	expectedSums := []float64{0, 0, 0, 0}
	assert.Equal(t, expectedSums, skewness.core.sums)
	assert.Equal(t, int(0), skewness.core.count)
	assert.Equal(t, uint64(0), skewness.core.queue.Len())
}

func TestSkewnessString(t *testing.T) {
	skewness := NewSkewness(3)
	expectedString := "moment.Skewness_{window:3}"
	assert.Equal(t, expectedString, skewness.String())
}
