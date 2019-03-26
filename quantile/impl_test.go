package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplInit(t *testing.T) {
	t.Run("pass: AVL implementation is supported", func(t *testing.T) {
		i := AVL
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("pass: red black implementation is supported", func(t *testing.T) {
		i := RB
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("pass: skip list implementation is supported", func(t *testing.T) {
		i := SkipList
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("fail: unsupported implementations return an error", func(t *testing.T) {
		i := Impl(-1)
		_, err := i.init()
		assert.EqualError(t, err, fmt.Sprintf("%v is not a supported Impl value", i))
	})
}
