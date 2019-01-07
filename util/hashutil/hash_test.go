package hashutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type item struct {
	hash uint64
	val  int
}

func (i item) Hash() uint64 {
	return i.hash
}

func (i item) Equal(m Mappable) bool {
	if i.Hash() != m.Hash() {
		return false
	}

	itemType := reflect.TypeOf(item{})
	mType := reflect.TypeOf(m)
	if itemType != mType {
		return false
	}

	mItem := m.(item)
	return i.val == mItem.val
}

func TestMapAdd(t *testing.T) {
	m := NewMap()
	item1 := item{hash: 1, val: 1}
	val1 := 17

	m.Add(item1, val1)

	assert.Equal(t, map[uint64][]*kv{
		item1.hash: []*kv{
			&kv{key: item1, value: val1},
		},
	}, m.hashmap)

	item2 := item{hash: 1, val: 2}
	val2 := 18
	m.Add(item2, val2)

	assert.Equal(t, map[uint64][]*kv{
		item1.hash: []*kv{
			&kv{key: item1, value: val1},
			&kv{key: item2, value: val2},
		},
	}, m.hashmap)
}

func TestMapContains(t *testing.T) {
	m := NewMap()
	i := item{hash: 1, val: 1}
	m.hashmap = map[uint64][]*kv{
		i.hash: []*kv{
			&kv{key: i, value: 17},
		},
	}

	assert.True(t, m.Contains(i))
	assert.False(t, m.Contains(item{hash: 2, val: 2}))
	assert.False(t, m.Contains(item{hash: i.hash, val: 2}))
}

func TestMapGet(t *testing.T) {
	m := NewMap()
	i := item{hash: 1, val: 1}
	m.hashmap = map[uint64][]*kv{
		i.hash: []*kv{
			&kv{key: i, value: 17},
		},
	}

	value, ok := m.Get(i)
	assert.True(t, ok)
	assert.Equal(t, 17, value)

	_, ok = m.Get(item{hash: 2, val: 2})
	assert.False(t, ok)

	_, ok = m.Get(item{hash: i.hash, val: 2})
	assert.False(t, ok)
}
