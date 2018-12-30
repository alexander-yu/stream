package ost

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplEmptyTree(t *testing.T) {
	t.Run("pass: AVL implementation is supported", func(t *testing.T) {
		i := AVL
		_, err := i.EmptyTree()
		assert.NoError(t, err)
	})

	t.Run("fail: unsupported implementations return an error", func(t *testing.T) {
		i := Impl(-1)
		_, err := i.EmptyTree()
		assert.EqualError(t, err, fmt.Sprintf("%v is not a supported OST implementation", i))
	})
}
