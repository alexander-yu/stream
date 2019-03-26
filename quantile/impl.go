package quantile

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream/quantile/order"
	"github.com/alexander-yu/stream/quantile/ost/avl"
	"github.com/alexander-yu/stream/quantile/ost/rb"
	"github.com/alexander-yu/stream/quantile/skiplist"
)

// Impl represents an enum that enumerates the currently supported implementations
// for the order.Statistic interface.
type Impl int

const (
	// AVL represents the AVL tree implementation for the order.Statistic interface
	AVL Impl = iota
	// RedBlack represents the red black tree implementation for the order.Statistic interface
	RedBlack
	// SkipList represents the skip list implementation for the order.Statistic interface
	SkipList
)

// Ptr returns a pointer to the Impl value.
func (i Impl) Ptr() *Impl {
	return &i
}

// Valid returns whether or not the Impl value is a valid value.
func (i Impl) Valid() bool {
	switch i {
	case AVL, RedBlack, SkipList:
		return true
	default:
		return false
	}
}

// Init returns an empty Impl struct, depending on which implementation
// is being called.
func (i Impl) init() (order.Statistic, error) {
	switch i {
	case AVL:
		return &avl.Tree{}, nil
	case RedBlack:
		return &rb.Tree{}, nil
	case SkipList:
		return skiplist.New()
	default:
		return nil, errors.Errorf("%v is not a supported Impl value", i)
	}
}
