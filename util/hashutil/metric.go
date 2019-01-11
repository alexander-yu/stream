package hashutil

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
)

// HashMetric computes a hash for a metric.
func HashMetric(name string, params map[string]interface{}) uint64 {
	paramStrings := sort.StringSlice{}
	for param, value := range params {
		paramStrings = append(paramStrings, fmt.Sprintf("%s:%v", param, value))
	}
	paramStrings.Sort()
	metricString := fmt.Sprintf("%s_%s", name, strings.Join(paramStrings, ","))
	h := fnv.New64a()
	h.Write([]byte(metricString))
	return h.Sum64()
}
