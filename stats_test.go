package stream

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var precision = 9

func testStats() *Stats {
	stats, _ := NewStats(&StatsConfig{
		Sums: map[int]bool{
			-1: true,
			0:  true,
			1:  true,
			2:  true,
			3:  true,
			4:  true,
		},
		Window: IntPtr(3),
		Median: BoolPtr(true),
	})

	for i := 1.; i < 5; i++ {
		stats.Push(i)
	}

	return stats
}

func roundFloat(x float64, n int) float64 {
	unit := 5 * math.Pow10(-n-1)
	return math.Round(x/unit) * unit
}

func approx(t *testing.T, x float64, y float64) {
	x = roundFloat(x, precision)
	y = roundFloat(y, precision)
	assert.Equal(t, x, y)
}

func TestPush(t *testing.T) {
	stats := testStats()
	median, _ := stats.medianStats.Median()

	sums := map[int]float64{
		-1: 13. / 12.,
		0:  3.,
		1:  9.,
		2:  29.,
		3:  99.,
		4:  353.,
	}

	assert.Equal(t, 4, stats.count)
	approx(t, 1., stats.min)
	approx(t, 4., stats.max)
	approx(t, 2.5, median)

	for k, sum := range sums {
		approx(t, sum, stats.sums[k])
	}
}

func TestCount(t *testing.T) {
	stats := testStats()
	assert.Equal(t, 4, stats.Count())
}

func TestMin(t *testing.T) {
	stats := testStats()
	approx(t, 1., stats.Min())
}

func TestMax(t *testing.T) {
	stats := testStats()
	approx(t, 4., stats.Max())
}

func TestSum(t *testing.T) {
	stats := testStats()
	sum, _ := stats.Sum(-1)
	approx(t, 13./12., sum)
}

func TestMean(t *testing.T) {
	stats := testStats()
	mean, _ := stats.Mean()
	approx(t, 3., mean)
}

func TestMoment(t *testing.T) {
	stats := testStats()
	stats.Push(8)
	moment, _ := stats.Moment(2)
	approx(t, 7., moment)
}

func TestStd(t *testing.T) {
	stats := testStats()
	stats.Push(8)
	std, _ := stats.Std()
	approx(t, math.Sqrt(7.), std)
}

func TestSkewness(t *testing.T) {
	stats := testStats()
	stats.Push(8)
	skewness, _ := stats.Skewness()

	adjust := 3.
	moment := 9.
	variance := 7.

	approx(t, adjust*moment/math.Pow(variance, 1.5), skewness)
}

func TestKurtosis(t *testing.T) {
	stats := testStats()
	stats.Push(8)
	kurtosis, _ := stats.Kurtosis()

	moment := 98. / 3.
	variance := 14. / 3.

	approx(t, moment/math.Pow(variance, 2.)-3., kurtosis)
}
