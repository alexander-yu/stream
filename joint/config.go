package joint

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// SumsConfig is an alias for a slice of Tuples; this configures
// the multinomial sums that a Core object will track.
type SumsConfig []Tuple

func (s1 *SumsConfig) add(s2 SumsConfig) {
	for _, t := range s2 {
		*s1 = append(*s1, t)
	}
}

// CoreConfig is the struct containing configuration options for
// instantiating a Core object.
type CoreConfig struct {
	Sums   SumsConfig
	Window *int
}

var defaultConfig = &CoreConfig{
	Sums:   SumsConfig{},
	Window: stream.IntPtr(0),
}

// MergeConfigs merges CoreConfig objects.
func MergeConfigs(configs ...*CoreConfig) (*CoreConfig, error) {
	switch len(configs) {
	case 0:
		return nil, errors.New("no configs available to merge")
	case 1:
		return configs[0], nil
	default:
		var window *int
		mergedConfig := &CoreConfig{
			Sums: SumsConfig{},
		}

		for _, config := range configs {
			if config.Sums != nil {
				mergedConfig.Sums.add(config.Sums)
			}

			if config.Window != nil {
				if window == nil {
					window = stream.IntPtr(*config.Window)
				} else if *window != *config.Window {
					return nil, errors.New("configs have differing windows")
				}
			}
		}

		mergedConfig.Window = window
		return mergedConfig, nil
	}
}

func validateConfig(config *CoreConfig) error {
	if config.Window != nil && *config.Window < 0 {
		return errors.Errorf("config has a negative window of %d", *config.Window)
	}

	for _, tuple := range config.Sums {
		for _, k := range tuple {
			if k < 0 {
				// The reason we allow for k = 0 here (even though there is no such
				// thing as a "0th moment") is because we can use it to ignore
				// variables; for example, we can have a Tuple of {2, 0, 0} represent a
				// calculation of the variance (or equivalently the sum of squared differences)
				// of the first variable, and simply ignore the other two. However, we still
				// need to make sure that an actual moment is being calculated, i.e. some
				// element of the Tuple is still positive.
				return errors.Errorf("config has a Tuple with a negative exponent of %d", k)
			}
		}

		if tuple.abs() == 0 {
			return errors.New("config has a Tuple that is all 0s (i.e. skips all variables)")
		}
	}

	return nil
}

func setConfigDefaults(config *CoreConfig) *CoreConfig {
	if config.Sums == nil {
		config.Sums = defaultConfig.Sums
	}

	if config.Window == nil {
		config.Window = defaultConfig.Window
	}

	return config
}
