package skiplist

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestNewSkipList(t *testing.T) {
	t.Run("fail: invalid Option is invalid", func(t *testing.T) {
		_, err := New(ProbabilityOption(-1))
		testutil.ContainsError(t, err, "error setting option")
	})

	t.Run("pass: no options is valid", func(t *testing.T) {
		_, err := New()
		require.NoError(t, err)
	})

	t.Run("pass: valid Options are valid", func(t *testing.T) {
		_, err := New(ProbabilityOption(0.25), MaxLevelOption(12))
		require.NoError(t, err)
	})
}

type SkipListSuite struct {
	suite.Suite
	skiplist *SkipList
}

func TestSkipListSuite(t *testing.T) {
	suite.Run(t, &SkipListSuite{})
}

func (s *SkipListSuite) SetupTest() {
	var err error
	s.skiplist, err = New(RandOption(rand.New(rand.NewSource(1))))
	s.NoError(err)
	s.skiplist.Add(5)
	s.skiplist.Add(6)
	s.skiplist.Add(7)
	s.skiplist.Add(3)
	s.skiplist.Add(4)
	s.skiplist.Add(1)
	s.skiplist.Add(2)
	s.skiplist.Add(1)
}

func (s *SkipListSuite) TestAdd() {
	s.Equal(8, s.skiplist.Size())
	s.Equal(
		strings.Join([]string{
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head---------tail",
			"head-1.000000000e+00--2.000000000e+00------tail",
			"head-1.000000000e+00-1.000000000e+00-2.000000000e+00-3.000000000e+00-4.000000000e+00-5.000000000e+00-6.000000000e+00-7.000000000e+00-tail",
			"",
		}, "\n"),
		s.skiplist.String(),
	)
	s.skiplist.Add(6.5)
	s.skiplist.Add(6.75)
	s.skiplist.Add(6.25)
	s.Equal(11, s.skiplist.Size())
	s.Equal(
		strings.Join([]string{
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head------------tail",
			"head-1.000000000e+00--2.000000000e+00------6.500000000e+00---tail",
			"head-1.000000000e+00-1.000000000e+00-2.000000000e+00-3.000000000e+00-4.000000000e+00-5.000000000e+00-6.000000000e+00-6.250000000e+00-6.500000000e+00-6.750000000e+00-7.000000000e+00-tail",
			"",
		}, "\n"),
		s.skiplist.String(),
	)
}

func (s *SkipListSuite) TestRemove() {
	s.Run("pass: successfully removes values", func() {
		s.SetupTest()
		s.skiplist.Add(6.5)
		s.skiplist.Add(6.75)
		s.skiplist.Add(6.25)

		s.skiplist.Remove(4)
		s.skiplist.Remove(1)
		s.skiplist.Remove(6.5)
		s.skiplist.Remove(6.75)
		s.skiplist.Remove(6.25)

		s.Equal(6, s.skiplist.Size())
		s.Equal(
			strings.Join([]string{
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head-------tail",
				"head--2.000000000e+00-----tail",
				"head-1.000000000e+00-2.000000000e+00-3.000000000e+00-5.000000000e+00-6.000000000e+00-7.000000000e+00-tail",
				"",
			}, "\n"),
			s.skiplist.String(),
		)

		s.skiplist.Remove(5)

		s.Equal(5, s.skiplist.Size())
		s.Equal(
			strings.Join([]string{
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head------tail",
				"head--2.000000000e+00----tail",
				"head-1.000000000e+00-2.000000000e+00-3.000000000e+00-6.000000000e+00-7.000000000e+00-tail",
				"",
			}, "\n"),
			s.skiplist.String(),
		)

		s.skiplist.Remove(6)
		s.skiplist.Remove(2)
		s.skiplist.Remove(3)
		s.skiplist.Remove(7)

		s.Equal(1, s.skiplist.Size())
		s.Equal(
			strings.Join([]string{
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head--tail",
				"head-1.000000000e+00-tail",
				"",
			}, "\n"),
			s.skiplist.String(),
		)

	})

	s.Run("pass: removing non-existent value is a no-op", func() {
		s.SetupTest()
		s.skiplist.Remove(8)

		s.Equal(8, s.skiplist.Size())
		s.Equal(
			strings.Join([]string{
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head---------tail",
				"head-1.000000000e+00--2.000000000e+00------tail",
				"head-1.000000000e+00-1.000000000e+00-2.000000000e+00-3.000000000e+00-4.000000000e+00-5.000000000e+00-6.000000000e+00-7.000000000e+00-tail",
				"",
			}, "\n"),
			s.skiplist.String(),
		)
	})
}

func (s *SkipListSuite) TestRank() {
	rank := s.skiplist.Rank(3)
	s.Equal(3, rank)

	rank = s.skiplist.Rank(5.5)
	s.Equal(6, rank)

	rank = s.skiplist.Rank(-1)
	s.Equal(0, rank)
}

func (s *SkipListSuite) TestSelect() {
	node := s.skiplist.Select(5)
	s.Equal(float64(5), node.Value())

	node = s.skiplist.Select(-1)
	s.Nil(node)

	node = s.skiplist.Select(9)
	s.Nil(node)
}

func (s *SkipListSuite) TestClear() {
	s.skiplist.Clear()
	s.Equal(0, s.skiplist.length)
	for _, node := range s.skiplist.prevs {
		s.Nil(node)
	}
	for _, node := range s.skiplist.head.next {
		s.Nil(node)
	}
	for _, width := range s.skiplist.head.width {
		s.Equal(1, width)
	}
}
