package moment

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewKurtosis(t *testing.T) {
	window := 3
	kurtosis := NewKurtosis(window)
	assert.Equal(t, &Kurtosis{
		variance: New(2, window),
		moment4:  New(4, window),
		config: &CoreConfig{
			Sums: SumsConfig{
				2: true,
				4: true,
			},
			Window: &window,
		},
	}, kurtosis)
}

type KurtosisPushSuite struct {
	suite.Suite
	kurtosis *Kurtosis
}

func TestKurtosisPushSuite(t *testing.T) {
	suite.Run(t, &KurtosisPushSuite{})
}

func (s *KurtosisPushSuite) SetupTest() {
	s.kurtosis = NewKurtosis(3)
	err := Init(s.kurtosis)
	s.Require().NoError(err)
}

func (s *KurtosisPushSuite) TestPushSuccess() {
	err := s.kurtosis.Push(3.)
	s.NoError(err)
}

func (s *KurtosisPushSuite) TestPushFailOnNullCore() {
	kurtosis := NewKurtosis(3)
	err := kurtosis.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *KurtosisPushSuite) TestPushFailOnQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.kurtosis.core.queue.Dispose()

	err := s.kurtosis.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *KurtosisPushSuite) TestPushFailOnQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.kurtosis.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.kurtosis.core.queue.Dispose()

	err := s.kurtosis.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

type KurtosisValueSuite struct {
	suite.Suite
	kurtosis *Kurtosis
}

func TestKurtosisValueSuite(t *testing.T) {
	suite.Run(t, &KurtosisValueSuite{})
}

func (s *KurtosisValueSuite) SetupTest() {
	s.kurtosis = NewKurtosis(3)
	err := Init(s.kurtosis)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.kurtosis.Push(x)
		s.Require().NoError(err)
	}
}

func (s *KurtosisValueSuite) TestValueSuccess() {
	value, err := s.kurtosis.Value()
	s.Require().NoError(err)

	moment := 98. / 3.
	variance := 14. / 3.

	testutil.Approx(s.T(), moment/math.Pow(variance, 2.)-3., value)
}

func (s *KurtosisValueSuite) TestValueFailOnNullCore() {
	kurtosis := NewKurtosis(3)
	_, err := kurtosis.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *KurtosisValueSuite) TestValueFailIfNoValuesSeen() {
	kurtosis := NewKurtosis(3)
	err := Init(kurtosis)
	s.Require().NoError(err)

	_, err = kurtosis.Value()
	testutil.ContainsError(s.T(), err, "no values seen yet")
}

func TestKurtosisClear(t *testing.T) {
	kurtosis := NewKurtosis(3)
	err := Init(kurtosis)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := kurtosis.Push(x)
		require.NoError(t, err)
	}

	kurtosis.Clear()
	expectedSums := []float64{0, 0, 0, 0, 0}
	assert.Equal(t, expectedSums, kurtosis.core.sums)
	assert.Equal(t, int(0), kurtosis.core.count)
	assert.Equal(t, uint64(0), kurtosis.core.queue.Len())
}

func TestKurtosisString(t *testing.T) {
	kurtosis := NewKurtosis(3)
	expectedString := "moment.Kurtosis_{window:3}"
	assert.Equal(t, expectedString, kurtosis.String())
}
