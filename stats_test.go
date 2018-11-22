package stream

import (
	"math"
	"testing"
)

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
