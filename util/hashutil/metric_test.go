package hashutil

import (
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMetric(t *testing.T) {
	name := "MockMetric"
	params := map[string]interface{}{
		"intParam":   int64(1),
		"floatParam": float64(3.5),
		"sliceParam": []int64{1, 2, 3},
	}
	expectedString := fmt.Sprintf(
		"%s_%s:%v,%s:%v,%s:%v",
		name,
		"floatParam",
		params["floatParam"],
		"intParam",
		params["intParam"],
		"sliceParam",
		params["sliceParam"],
	)
	h := fnv.New64a()
	h.Write([]byte(expectedString))
	expectedHash := h.Sum64()

	hash := HashMetric(name, params)
	assert.Equal(t, expectedHash, hash)
}
