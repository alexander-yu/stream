package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewAutocov(t *testing.T) {
	t.Run("pass: valid Autocov is valid", func(t *testing.T) {
		autocov, err := NewAutocov(1, 3)
		require.NoError(t, err)
		assert.Equal(t, 1, autocov.lag)
		assert.Equal(t, NewCov(3), autocov.cov)
	})

	t.Run("fail: negative lag returns error", func(t *testing.T) {
		_, err := NewAutocov(-1, 3)
		testutil.ContainsError(t, err, "-1 is a negative lag")
	})
}

func TestNewGlobalAutocov(t *testing.T) {
	autocov, err := NewAutocov(1, 0)
	require.NoError(t, err)
	globalAutocov, err := NewGlobalAutocov(1)
	require.NoError(t, err)

	assert.Equal(t, autocov, globalAutocov)
}

type AutocovPushSuite struct {
	suite.Suite
	autocov  *Autocov
	autocov0 *Autocov
}

func TestAutocovPushSuite(t *testing.T) {
	suite.Run(t, &AutocovPushSuite{})
}

func (s *AutocovPushSuite) SetupTest() {
	var err error
	s.autocov, err = NewAutocov(1, 3)
	s.Require().NoError(err)
	err = Init(s.autocov)
	s.Require().NoError(err)

	s.autocov0, err = NewAutocov(0, 3)
	s.Require().NoError(err)
	err = Init(s.autocov0)
	s.Require().NoError(err)
}

func (s *AutocovPushSuite) TestPushSuccess() {
	err := s.autocov.Push(3.)
	s.NoError(err)
}

func (s *AutocovPushSuite) TestPushFailOnNullCore() {
	autocov, err := NewAutocov(1, 3)
	s.Require().NoError(err)
	err = autocov.Push(0.)
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *AutocovPushSuite) TestPushFailOnCorePushFailureForLag0() {
	// dispose the queue to simulate an error when we try to push to the core
	s.autocov0.core.queue.Dispose()
	err := s.autocov0.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocovPushSuite) TestPushFailOnCoreQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.autocov.core.queue.Dispose()

	// no error yet because we have not populated the lag yet
	err := s.autocov.Push(8.)
	s.Require().NoError(err)

	err = s.autocov.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocovPushSuite) TestPushFailOnCoreQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.autocov.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.autocov.core.queue.Dispose()

	err := s.autocov.Push(3.)
	testutil.ContainsError(s.T(), err, "error pushing to core")
}

func (s *AutocovPushSuite) TestPushFailOnLagQueueInsertionFailure() {
	// dispose the queue to simulate an error when we try to insert into the queue
	s.autocov.queue.Dispose()

	val := 8.
	err := s.autocov.Push(val)
	testutil.ContainsError(s.T(), err, fmt.Sprintf("error pushing %f to lag queue", val))
}

func (s *AutocovPushSuite) TestPushFailOnLagQueueRetrievalFailure() {
	xs := []float64{1, 2, 3}
	for _, x := range xs {
		err := s.autocov.Push(x)
		s.Require().NoError(err)
	}

	// dispose the queue to simulate an error when we try to retrieve from the queue
	s.autocov.queue.Dispose()

	err := s.autocov.Push(3.)
	testutil.ContainsError(s.T(), err, "error popping item from lag queue")
}

type AutocovValueSuite struct {
	suite.Suite
	autocov  *Autocov
	autocov0 *Autocov
}

func TestAutocovValueSuite(t *testing.T) {
	suite.Run(t, &AutocovValueSuite{})
}

func (s *AutocovValueSuite) SetupTest() {
	var err error
	s.autocov, err = NewAutocov(1, 3)
	s.Require().NoError(err)
	err = Init(s.autocov)
	s.Require().NoError(err)

	s.autocov0, err = NewAutocov(0, 3)
	s.Require().NoError(err)
	err = Init(s.autocov0)
	s.Require().NoError(err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := s.autocov.Push(x)
		s.Require().NoError(err)

		err = s.autocov0.Push(x)
		s.Require().NoError(err)
	}
}

func (s *AutocovValueSuite) TestValueSuccess() {
	value, err := s.autocov.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 5./2., value)
}

func (s *AutocovValueSuite) TestValueSuccessForLag0() {
	value, err := s.autocov0.Value()
	s.Require().NoError(err)
	testutil.Approx(s.T(), 7., value)
}

func (s *AutocovValueSuite) TestValueFailOnNullCore() {
	autocov, err := NewAutocov(1, 3)
	s.Require().NoError(err)
	_, err = autocov.Value()
	testutil.ContainsError(s.T(), err, "Core is not set")
}

func (s *AutocovValueSuite) TestValueFailIfNotEnoughValuesSeen() {
	autocov, err := NewAutocov(1, 3)
	s.Require().NoError(err)
	err = Init(autocov)
	s.Require().NoError(err)

	err = autocov.Push(1)
	s.Require().NoError(err)

	_, err = autocov.Value()
	testutil.ContainsError(s.T(), err, fmt.Sprintf(
		"Not enough values seen; at least %d observations must be made",
		2,
	))
}

func TestAutocovClear(t *testing.T) {
	autocov, err := NewAutocov(1, 3)
	require.NoError(t, err)
	err = Init(autocov)
	require.NoError(t, err)

	xs := []float64{1, 2, 3, 4, 8}
	for _, x := range xs {
		err := autocov.Push(x)
		require.NoError(t, err)
	}

	autocov.Clear()

	expectedSums := map[uint64]float64{
		0:  0.,
		1:  0.,
		31: 0.,
		32: 0.,
	}
	assert.Equal(t, expectedSums, autocov.core.sums)
	assert.Equal(t, expectedSums, autocov.core.newSums)
	assert.Equal(t, 0, autocov.core.count)
	assert.Equal(t, uint64(0), autocov.core.queue.Len())
	assert.Equal(t, uint64(0), autocov.queue.Len())
}

func TestAutocovString(t *testing.T) {
	autocov, err := NewAutocov(1, 3)
	require.NoError(t, err)
	expectedString := "joint.Autocov_{lag:1,window:3}"
	assert.Equal(t, expectedString, autocov.String())
}
