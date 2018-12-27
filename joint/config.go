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
	Vars   *int
}

var defaultConfig = &CoreConfig{
	Sums:   SumsConfig{},
	Window: stream.IntPtr(0),
	Vars:   nil,
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
		var vars *int
		mergedConfig := &CoreConfig{
			Sums: SumsConfig{},
		}

		for _, config := range configs {
			if config.Sums != nil {
				mergedConfig.Sums.add(config.Sums)
			}

			if config.Window != nil {
				if window == nil {
					window = config.Window
				} else if *window != *config.Window {
					return nil, errors.New("configs have differing windows")
				}
			}

			if config.Vars != nil {
				if vars == nil {
					vars = config.Vars
				} else if *vars != *config.Vars {
					return nil, errors.New("configs have differing vars")
				}
			}
		}

		// remove any duplicate Tuples, or Tuples that are less than or equal
		// to other Tuples, since they'll get tracked automatically
		tupleMap := map[int]bool{}
		sums := SumsConfig{}
		for i, m := range mergedConfig.Sums {
			// remove dupes
			if _, ok := tupleMap[m.hash()]; ok {
				continue
			}

			// if Tuple is less than or equal to an already existing Tuple,
			// skip this one
			leq := false
			for j, n := range mergedConfig.Sums {
				if i <= j {
					continue
				}

				allLeq := true
				for i := range m {
					if m[i] > n[i] {
						allLeq = false
						break
					}
				}
				if allLeq {
					leq = true
					break
				}
			}
			if leq {
				continue
			}

			tupleMap[m.hash()] = true
			sums = append(sums, m)
		}

		mergedConfig.Sums = sums
		mergedConfig.Window = window
		mergedConfig.Vars = vars
		return mergedConfig, nil
	}
}

func validateConfig(config *CoreConfig) error {
	if config.Window == nil {
		return errors.New("config Window is not set")
	} else if *config.Window < 0 {
		return errors.Errorf("config has a negative window of %d", *config.Window)
	}

	if config.Vars == nil {
		return errors.New("config Vars is not set")
	} else if *config.Vars < 2 {
		return errors.Errorf("config has less than 2 vars: %d < 2", *config.Vars)
	}

	for _, tuple := range config.Sums {
		if len(tuple) != *config.Vars {
			return errors.Errorf(
				"config has a Tuple (%v) with length %d but Vars = %d",
				tuple,
				len(tuple),
				*config.Vars,
			)
		}

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

	if config.Vars == nil {
		config.Vars = defaultConfig.Vars
	}

	return config
}
