package skiplist

import (
	"math/rand"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream/quantile/order"
)

// MaxLevelOption creates an option that sets the max level for a skiplist.
func MaxLevelOption(maxLevel int) order.Option {
	return func(s order.Statistic) error {
		var (
			skiplist *SkipList
			ok       bool
		)
		if skiplist, ok = s.(*SkipList); !ok {
			return errors.New("attempted to set max level on a non-skiplist")
		} else if maxLevel < 1 || maxLevel > 64 {
			return errors.Errorf("attempted to set max level %d not in [1, 64]", maxLevel)
		}
		skiplist.maxLevel = maxLevel
		return nil
	}
}

// ProbabilityOption creates an option that sets the probability for deciding
// levels in a skiplist.
func ProbabilityOption(p float64) order.Option {
	return func(s order.Statistic) error {
		var (
			skiplist *SkipList
			ok       bool
		)
		if skiplist, ok = s.(*SkipList); !ok {
			return errors.New("attempted to set probability on a non-skiplist")
		} else if p <= 0 || p >= 1 {
			return errors.Errorf("attempted to set probability %f not in (0, 1)", p)
		}
		skiplist.p = p
		return nil
	}
}

// RandOption creates an option that sets the rand source for the skip list.
func RandOption(r *rand.Rand) order.Option {
	return func(s order.Statistic) error {
		var (
			skiplist *SkipList
			ok       bool
		)
		if skiplist, ok = s.(*SkipList); !ok {
			return errors.New("attempted to set rand source on a non-skiplist")
		}
		skiplist.rand = r
		return nil
	}
}
