package skiplist

import "math/rand"

// Option is an optional argument to New,
// which sets an optional field for creating a SkipList
type Option func(*SkipList)

// MaxLevelOption creates an option that sets the max level for a skiplist.
func MaxLevelOption(maxLevel int) Option {
	return func(s *SkipList) {
		s.maxLevel = maxLevel
	}
}

// ProbabilityOption creates an option that sets the probability for deciding
// levels in a skiplist.
func ProbabilityOption(p float64) Option {
	return func(s *SkipList) {
		s.p = p
	}
}

// RandOption creates an option that sets the rand source for the skip list.
func RandOption(r *rand.Rand) Option {
	return func(s *SkipList) {
		s.rand = r
	}
}
