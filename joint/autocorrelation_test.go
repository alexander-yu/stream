package joint

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewAutocorr(t *testing.T) {
	t.Run("pass: valid Autocorr is valid", func(t *testing.T) {
		autocorr, err := NewAutocorr(1, 3)
		require.NoError(t, err)
		assert.Equal(t, 1, autocorr.lag)
		assert.Equal(t, NewCorr(3), autocorr.corr)
	})

	t.Run("fail: negative lag returns error", func(t *testing.T) {
		_, err := NewAutocorr(-1, 3)
		testutil.ContainsError(t, err, "-1 is a negative lag")
	})
}

type AutocorrPushSuite struct {
	suite.Suite
	autocorr  *Autocorr
	autocorr0 *Autocorr
}

func TestAutocorrPushSuite(t *testing.T) {
	suite.Run(t, &AutocorrPushSuite{})
}

func (s *AutocorrPushSuite) SetupTest() {
	var err error
	s.autocorr, err = NewAutocorr(1, 3)
	s.Require().NoError(err)
	err = Init(s.autocorr)
	s.Require().NoError(err)

	s.autocorr0, err = NewAutocorr(0, 3)
	s.Require().NoError(err)
	err = Init(s.autocorr0)
	s.Require().NoError(err)
}

func (s *AutocorrPushSuite) TestPushSuccess() {
	err := s.autocorr.Push(3.)
	s.NoError(err)
}

func (s *AutocorrPushSuite) TestPushFailOnNullCore() {
	autocorr, err := NewAutocorr(1, 3)
	s.Require().NoError(err)
	err = autocorr.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *AutocorrPushSuite) TestPushFailOnCorePushFailureForLag0() {
	// dispose the queue to simulate an error when we try to push to the core
	s.autocorr0.core.queue.Dispose()
	err := s.autocorr0.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocorrPushSuite) TestPushFailOnCoreQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.autocorr.core.queue.Dispose()

	// no error yet because we have not populated the lag yet
	err := s.autocorr.Push(8.)
	s.Require().NoError(err)

	err = s.autocorr.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocorrPushSuite) TestPushFailOnCoreQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.autocorr.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.autocorr.core.queue.Dispose()

	err := s.autocorr.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocorrPushSuite) TestPushFailOnLagQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.autocorr.queue.Dispose()

	val := 8.
	err := s.autocorr.Push(val)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %f to lag queue", val))
}

func (s *AutocorrPushSuite) TestPushFailOnLagQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.autocorr.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.autocorr.queue.Dispose()

	err := s.autocorr.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from lag queue")
}

type AutocorrValueSuite struct {
	suite.Suite
	autocorr  *Autocorr
	autocorr0 *Autocorr
}

func TestAutocorrValueSuite(t *testing.T) {
	suite.Run(t, &AutocorrValueSuite{})
}

func (s *AutocorrValueSuite) SetupTest() {
	var err error
	s.autocorr, err = NewAutocorr(1, 3)
	s.Require().NoError(err)
	err = Init(s.autocorr)
	s.Require().NoError(err)

	s.autocorr0, err = NewAutocorr(0, 3)
	s.Require().NoError(err)
	err = Init(s.autocorr0)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.autocorr.Push(x)
		s.Require().NoError(err)

		err = s.autocorr0.Push(x)
		s.Require().NoError(err)
	}
}

func (s *AutocorrValueSuite) TestValueSuccess() {
	value, err := s.autocorr.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 5./(2.*math.Sqrt(7.)), value)
}

func (s *AutocorrValueSuite) TestValueSuccessForLag0() {
	value, err := s.autocorr0.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 1., value)
}

func (s *AutocorrValueSuite) TestValueFailOnNullCore() {
	autocorr, err := NewAutocorr(1, 3)
	s.Require().NoError(err)
	_, err = autocorr.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *AutocorrValueSuite) TestValueFailIfNotEnoughValuesSeen() {
	autocorr, err := NewAutocorr(1, 3)
	s.Require().NoError(err)
	err = Init(autocorr)
	s.Require().NoError(err)

	err = autocorr.Push(1)
	s.Require().NoError(err)

	_, err = autocorr.Value()
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Not enough values seen; at least %d observations must be made",
		2,
	))
}

func TestAutocorrClear(t *testing.T) {
	autocorr, err := NewAutocorr(1, 3)
	require.NoError(t, err)
	err = Init(autocorr)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := autocorr.Push(x)
		require.NoError(t, err)
	}

	autocorr.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		2:  0.,
		31: 0.,
		32: 0.,
		62: 0.,
	}
	assert.Equal(t, expectedSums, autocorr.core.sums)
	assert.Equal(t, expectedSums, autocorr.core.newSums)
	assert.Equal(t, 0, autocorr.core.count)
	assert.Equal(t, uint64(0), autocorr.core.queue.Len())
	assert.Equal(t, uint64(0), autocorr.queue.Len())
}

func TestAutocorrString(t *testing.T) {
	autocorr, err := NewAutocorr(1, 3)
	require.NoError(t, err)
	expectedString := "joint.Autocorr_{lag:1,window:3}"
	assert.Equal(t, expectedString, autocorr.String())
}
