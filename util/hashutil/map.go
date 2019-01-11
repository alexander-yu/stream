package hashutil

// Map implements a very basic hashmap with separate chaining.
// While go maps already implement separate chaining under the hood,
// maps only support a specific number of types, in particular those
// that are supported by the == and != operators. This allows for
// the mapping of custom types that satisfy the Mappable interface.
type Map struct {
	hashmap map[uint64][]*kv
}

// Mappable is the interface for items that can be stored in Map;
// it needs a hash method to store it into buckets, and an equality
// method to compare it with other elements in the same bucket.
type Mappable interface {
	Hash() uint64
	Equal(Mappable) bool
}

type kv struct {
	key   Mappable
	value interface{}
}

// NewMap initializes a new Map.
func NewMap() *Map {
	return &Map{hashmap: map[uint64][]*kv{}}
}

// Add inserts a key and value into the map.
func (m *Map) Add(key Mappable, value interface{}) {
	hash := key.Hash()
	item := &kv{key: key, value: value}

	if _, ok := m.hashmap[hash]; !ok {
		m.hashmap[hash] = []*kv{item}
		return
	}
	if !m.Contains(key) {
		m.hashmap[hash] = append(m.hashmap[hash], item)
	}
}

// Contains returns whether or not the key is contained in the map.
func (m *Map) Contains(key Mappable) bool {
	hash := key.Hash()
	bucket, ok := m.hashmap[hash]

	if !ok {
		return false
	}
	for _, kv := range bucket {
		if kv.key.Equal(key) {
			return true
		}
	}

	return false
}

// Get retrieves the value for a key in the map; returns nil and false if
// the key does not exist in the map.
func (m *Map) Get(key Mappable) (interface{}, bool) {
	hash := key.Hash()
	bucket, ok := m.hashmap[hash]

	if !ok {
		return nil, false
	}
	for _, kv := range bucket {
		if kv.key.Equal(key) {
			return kv.value, true
		}
	}

	return nil, false
}
