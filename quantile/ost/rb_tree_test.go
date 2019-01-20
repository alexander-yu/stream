package ost

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RBTreeSuite struct {
	suite.Suite
	tree *RBTree
}

func TestRBTreeSuite(t *testing.T) {
	suite.Run(t, &RBTreeSuite{})
}

func (s *RBTreeSuite) SetupTest() {
	s.tree = &RBTree{}
	s.tree.Add(5)
	s.tree.Add(6)
	s.tree.Add(7)
	s.tree.Add(3)
	s.tree.Add(4)
	s.tree.Add(1)
	s.tree.Add(2)
	s.tree.Add(1)
}

func (s *RBTreeSuite) TestRBTreeAdd() {
	s.Equal(8, s.tree.Size())
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
	s.Equal(
		strings.Join([]string{
			"│       ┌── 7.000000",
			"│   ┌── 6.750000",
			"│   │   │   ┌── 6.500000",
			"│   │   │   │   └── 6.250000",
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

func (s *RBTreeSuite) TestRBTreeRemove() {
	s.Run("pass: successfully removes values", func() {
		s.SetupTest()
		s.tree.Add(6.5)
		s.tree.Add(6.75)
		s.tree.Add(6.25)

		s.tree.Remove(4)
		s.tree.Remove(1)
		s.tree.Remove(6.5)
		s.tree.Remove(6.75)
		s.tree.Remove(6.25)

		s.Equal(6, s.tree.Size())
		s.Equal(
			strings.Join([]string{
				"│   ┌── 7.000000",
				"│   │   └── 6.000000",
				"└── 5.000000",
				"    │   ┌── 3.000000",
				"    └── 2.000000",
				"        └── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)

		s.tree.Remove(5)

		s.Equal(5, s.tree.Size())
		s.Equal(
			strings.Join([]string{
				"│   ┌── 7.000000",
				"└── 6.000000",
				"    │   ┌── 3.000000",
				"    └── 2.000000",
				"        └── 1.000000",
				"",
			}, "\n"),
			s.tree.String(),
		)

		s.tree.Remove(6)
		s.tree.Remove(2)
		s.tree.Remove(3)
		s.tree.Remove(7)

		s.Equal(1, s.tree.Size())
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

func (s *RBTreeSuite) TestRBTreeRank() {
	rank := s.tree.Rank(3)
	s.Equal(3, rank)

	rank = s.tree.Rank(5.5)
	s.Equal(6, rank)

	rank = s.tree.Rank(-1)
	s.Equal(0, rank)
}

func (s *RBTreeSuite) TestRBTreeSelect() {
	node := s.tree.Select(5)
	s.Equal(float64(5), node.Value())

	node = s.tree.Select(-1)
	s.Nil(node)

	node = s.tree.Select(9)
	s.Nil(node)
}

func (s *RBTreeSuite) TestRBTreeClear() {
	s.tree.Clear()
	s.Equal(&RBTree{}, s.tree)
}
