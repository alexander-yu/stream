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

// StatsConfig is the struct containing configuration options for
// instantiating a Stats object.
type StatsConfig struct {
	sums   map[int]bool
	window int
}

// NewStats returns a Stats object with initiated sums.
func NewStats(config *StatsConfig) (*Stats, error) {
	if len(config.sums) == 0 {
		return nil, errors.New("stream: map is empty")
	} else if config.window <= 0 {
		return nil, errors.New("stream: window size is nonpositive")
	}

	stats := Stats{min: math.Inf(1), max: math.Inf(-1)}
	stats.window = config.window
	stats.sums = make(map[int]float64)

	for k := range config.sums {
		stats.sums[k] = 0
	}

	return &stats, nil
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

	mean, err := s.Mean()
	if err != nil {
		return 0, err
	}

	var moment float64
	for i := 0; i <= k; i++ {
		sum, err := s.Sum(i)
		if err != nil {
			return 0, fmt.Errorf("stream: %d is not a tracked power sum", i)
		}

		moment += float64(binom(k, i)*sign(k-i)) * math.Pow(mean, float64(k-i)) * sum
	}

	// Ignore the error; if execution gets here, then s.Sum(0) should already have a valid result
	count, _ := s.Sum(0)
	moment /= count

	return moment, nil
}

// Std returns the running standard deviation of values seen.
func (s *Stats) Std() (float64, error) {
	variance, err := s.Moment(2)
	return math.Sqrt(variance), err
}

// Mean returns the running mean of values seen.
func (s *Stats) Mean() (float64, error) {
	count, err0 := s.Sum(0)
	sum, err1 := s.Sum(1)
	if err0 != nil {
		return 0, errors.New("stream: 0 is not a tracked power sum")
	} else if err1 != nil {
		return 0, errors.New("stream: 1 is not a tracked power sum")
	}

	mean := sum / count

	return mean, nil
}

// Skewness returns the skewness of values seen.
func (s *Stats) Skewness() (float64, error) {
	variance, err2 := s.Moment(2)
	moment, err3 := s.Moment(3)
	if err2 != nil {
		return 0, err2
	} else if err3 != nil {
		return 0, err3
	}

	return moment / math.Pow(variance, float64(3)/float64(2)), nil
}

// Kurtosis returns the kurtosis of values seen.
func (s *Stats) Kurtosis() (float64, error) {
	variance, err2 := s.Moment(2)
	moment, err4 := s.Moment(4)
	if err2 != nil {
		return 0, err2
	} else if err4 != nil {
		return 0, err4
	}

	return moment / math.Pow(variance, float64(2)), nil
}

// Clear clears all stats being tracked.
func (s *Stats) Clear() {
	for k := range s.sums {
		s.sums[k] = 0
	}

	s.count = 0
	s.min = 0
	s.max = 0
	s.vals = nil
}
