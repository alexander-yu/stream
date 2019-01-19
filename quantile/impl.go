package quantile

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream/quantile/order"
	"github.com/alexander-yu/stream/quantile/ost"
)

// Impl represents an enum that enumerates the currently supported implementations
// for the order.Statistic interface.
type Impl int

const (
	// AVL represents the AVL tree implementation for the Tree interface
	AVL Impl = iota
	// RB represents the red black tree implementation for the Tree interface
	RB
)

// Init returns an empty Impl struct, depending on which implementation
// is being called.
func (i Impl) init() (order.Statistic, error) {
	switch i {
	case AVL:
		return &ost.AVLTree{}, nil
	case RB:
		return &ost.RBTree{}, nil
	default:
		return nil, errors.Errorf("%v is not a supported OST implementation", i)
	}
}
