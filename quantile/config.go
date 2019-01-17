package quantile

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Config is the struct containing configuration options
// for quantile metrics.
type Config struct {
	Window        *int
	Interpolation *Interpolation
}

var defaultConfig = &Config{
	Window:        stream.IntPtr(0),
	Interpolation: Linear.Ptr(),
}

// Interpolation represents an enum that enumerates the
// different interpolation methods that can be chosen
// when retrieving quantile metrics. In particular,
// if φ is the quantile value and n is the number of elements,
// then the raw estimated index for the φ-quantile element
// is i' = φ * (n - 1). However, this is not guaranteed to be an integer
// value. Thus, if the quantile of a window of elements actually
// lies in between two elements, this determines how
// the returned value will be calculated (depending
// on those two elements).
type Interpolation int

const (
	// Linear performs linear interpolation between the two elements.
	// If the φ-quantile i' lies between i and i + 1, then we return
	// (1 - (i' - i)) * a_i + (i' - i) * a_(i + 1).
	Linear Interpolation = 1
	// Lower chooses the lower of the two elements.
	Lower Interpolation = 2
	// Higher chooses the higher of the two elements.
	Higher Interpolation = 3
	// Nearest chooses the element whose index is closest to i'. If
	// i' is the midpoint between i and i + 1, then ties are broken
	// based on the parity of i; if i is even, then we choose a_i, and
	// otherwise choose a_(i + 1).
	Nearest Interpolation = 4
	// Midpoint performs the average of the two elements.
	Midpoint Interpolation = 5
)

// Valid returns whether or not the Interpolation value is a valid
// value.
func (i Interpolation) Valid() bool {
	switch i {
	case Linear, Lower, Higher, Nearest, Midpoint:
		return true
	default:
		return false
	}
}

// Ptr returns a pointer to the Interpolation value.
func (i Interpolation) Ptr() *Interpolation {
	return &i
}

func validateConfig(config *Config) error {
	if config.Window == nil {
		return errors.New("config Window is not set")
	} else if config.Interpolation == nil {
		return errors.New("config Interpolation is not set")
	}

	if *config.Window < 0 {
		return errors.Errorf("config has a negative window of %d", *config.Window)
	}

	if !config.Interpolation.Valid() {
		return errors.Errorf(
			"config has an invalid Interpolation of %d",
			*config.Interpolation,
		)
	}

	return nil
}

func setConfigDefaults(config *Config) *Config {
	if config.Window == nil {
		config.Window = defaultConfig.Window
	}

	if config.Interpolation == nil {
		config.Interpolation = defaultConfig.Interpolation
	}

	return config
}
