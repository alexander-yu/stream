package stream

import (
	"errors"
	"fmt"
	"math"
)

// Stats is the struct that provides stats being tracked.
type Stats struct {
	sums   map[int]float64
	count  int
	min    float64
	max    float64
	window int
	vals   []float64
}

// NewStats returns a Stats object with initiated sums.
func NewStats(sums map[int]float64) (*Stats, error) {
	if len(sums) == 0 {
		return nil, errors.New("stream: map is empty")
	}

	return &Stats{sums: sums, min: math.Inf(1), max: math.Inf(-1)}, nil
}

// NewWindowedStats returns a Stats object with initiated sums and window size.
// Using a window will mean that power sums will only be calculated over the current
// running window; count/min/max will still be global values (i.e. over all values seen).
func NewWindowedStats(sums map[int]float64, window int) (*Stats, error) {
	if window <= 0 {
		return nil, errors.New("stream: window size is nonpositive")
	}

	stats, err := NewStats(sums)
	if err != nil {
		return nil, err
}

	stats.window = window
	return stats, nil
}

// Push adds a new value for a Stats object to consume.
func (s *Stats) Push(x float64) {
	if s.window != 0 {
		s.vals = append(s.vals, x)

		if len(s.vals) > s.window {
			tail := s.vals[0]
			s.vals = s.vals[1:]

			for k := range s.sums {
				s.sums[k] -= math.Pow(tail, float64(k))
			}
		}
	}

	for k := range s.sums {
		s.sums[k] += math.Pow(x, float64(k))
	}

	s.count++
	s.min = math.Min(s.min, x)
	s.max = math.Max(s.max, x)
}

// Count returns the number of values seen.
func (s *Stats) Count() int {
	return s.count
}

// Min returns the min of values seen.
func (s *Stats) Min() float64 {
	return s.min
}

// Max returns the max of values seen.
func (s *Stats) Max() float64 {
	return s.max
}

// Sum returns the running kth-power sum of values seen.
func (s *Stats) Sum(k int) (float64, error) {
	if sum, ok := s.sums[k]; ok {
		return sum, nil
	}

	return 0, fmt.Errorf("stream: %d is not a tracked power sum", k)
}

// Moment returns the running kth-moment of values seen.
func (s *Stats) Moment(k int) (float64, error) {
	if k < 0 {
		return 0, errors.New("stream: negative moment")
	} else if k == 0 {
		return s.Sum(0)
	}

	count, err0 := s.Sum(0)
	sum, err1 := s.Sum(1)
	if err0 != nil {
		return 0, errors.New("stream: 0 is not a tracked power sum")
	} else if err1 != nil {
		return 0, errors.New("stream: 1 is not a tracked power sum")
	}

	mean := sum / count

	var moment float64
	for i := 0; i <= k; i++ {
		sum, err := s.Sum(i)
		if err != nil {
			return 0, fmt.Errorf("stream: %d is not a tracked power sum", i)
		}

		moment += float64(binom(k, i)*sign(k-i)) * math.Pow(mean, float64(k-i)) * sum
	}

	moment /= count

	return moment, nil
}
