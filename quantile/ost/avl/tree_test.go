package avl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TreeSuite struct {
	suite.Suite
	tree *Tree
}

func TestTreeSuite(t *testing.T) {
	suite.Run(t, &TreeSuite{})
}

func (s *TreeSuite) SetupTest() {
	s.tree = &Tree{}
	s.tree.Add(5)
	s.tree.Add(6)
	s.tree.Add(7)
	s.tree.Add(3)
	s.tree.Add(4)
	s.tree.Add(1)
	s.tree.Add(2)
	s.tree.Add(1)
}

func (s *TreeSuite) TestAdd() {
	s.Equal(8, s.tree.Size())
	s.Equal(3, s.tree.Height())
	s.Equal(
		strings.Join([]string{
			"│       ┌── 7.000000",
			"│   ┌── 6.000000",
			"│   │   └── 5.000000",
			"└── 4.000000",
			"    │   ┌── 3.000000",
			"    └── 2.000000",
			"        └── 1.000000",
			"            └── 1.000000",
			"",
		}, "\n"),
		s.tree.String(),
	)

	s.tree.Add(6.5)
	s.tree.Add(6.75)
	s.tree.Add(6.25)
	s.Equal(11, s.tree.Size())
	s.Equal(3, s.tree.Height())
	s.Equal(
		strings.Join([]string{
			"│           ┌── 7.000000",
			"│       ┌── 6.750000",
			"│   ┌── 6.500000",
			"│   │   │   ┌── 6.250000",
			"│   │   └── 6.000000",
			"│   │       └── 5.000000",
			"└── 4.000000",
			"    │   ┌── 3.000000",
			"    └── 2.000000",
			"        └── 1.000000",
			"            └── 1.000000",
			"",
		}, "\n"),
		s.tree.String(),
	)
}

func (s *TreeSuite) TestRemove() {
	s.Run("pass: successfully removes values", func() {
		s.SetupTest()
		s.tree.Remove(5)
		s.tree.Remove(7)

		s.Equal(6, s.tree.Size())
		s.Equal(2, s.tree.Height())
		s.Equal(
			strings.Join([]string{
				"│       ┌── 6.000000",
				"│   ┌── 4.000000",
				"│   │   └── 3.000000",
				"└── 2.000000",
				"    └── 1.000000",
				"        └── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)

		s.tree.Remove(2)
		s.Equal(5, s.tree.Size())
		s.Equal(2, s.tree.Height())
		s.Equal(
			strings.Join([]string{
				"│       ┌── 6.000000",
				"│   ┌── 4.000000",
				"└── 3.000000",
				"    └── 1.000000",
				"        └── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)

		s.tree.Remove(1)
		s.tree.Remove(6)
		s.tree.Remove(4)
		s.tree.Remove(3)
		s.Equal(1, s.tree.Size())
		s.Equal(0, s.tree.Height())
		s.Equal(
			strings.Join([]string{
				"└── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)
	})

	s.Run("pass: removing non-existent value is a no-op", func() {
		s.SetupTest()
		s.tree.Remove(8)

		s.Equal(8, s.tree.Size())
		s.Equal(3, s.tree.Height())
		s.Equal(
			strings.Join([]string{
				"│       ┌── 7.000000",
				"│   ┌── 6.000000",
				"│   │   └── 5.000000",
				"└── 4.000000",
				"    │   ┌── 3.000000",
				"    └── 2.000000",
				"        └── 1.000000",
				"            └── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)
	})
}

func (s *TreeSuite) TestRank() {
	rank := s.tree.Rank(3)
	s.Equal(3, rank)

	rank = s.tree.Rank(5.5)
	s.Equal(6, rank)

	rank = s.tree.Rank(-1)
	s.Equal(0, rank)
}

func (s *TreeSuite) TestSelect() {
	node := s.tree.Select(5)
	s.Equal(float64(5), node.Value())

	node = s.tree.Select(-1)
	s.Nil(node)

	node = s.tree.Select(9)
	s.Nil(node)
}

func (s *TreeSuite) TestClear() {
	s.tree.Clear()
	s.Equal(&Tree{}, s.tree)
}
