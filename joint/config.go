package joint

import (
	"github.com/pkg/errors"
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
	Sums   SumsConfig // sums tracked must be positive, and must track > 1 variables
	Window *int       // must be 0 if decay is set, must be nonnegative in general
	Vars   *int       // must be inferrable from Sums if not set; otherwise must be > 1
	Decay  *float64   // optional, must lie in the interval (0, 1)
}

var defaultConfig = &CoreConfig{
	Sums:   SumsConfig{},
	Window: nil,
	Vars:   nil,
	Decay:  nil,
}

// MergeConfigs merges CoreConfig objects.
func MergeConfigs(configs ...*CoreConfig) (*CoreConfig, error) {
	switch len(configs) {
	case 0:
		return nil, errors.New("no configs available to merge")
	case 1:
		return configs[0], nil
	default:
		var (
			window *int
			vars   *int
			decay  *float64
		)
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

			if config.Decay != nil {
				if decay == nil {
					decay = config.Decay
				} else if *decay != *config.Decay {
					return nil, errors.New("configs have differing decays")
				}
			}
		}

		mergedConfig.Sums = simplifySums(mergedConfig.Sums)
		mergedConfig.Window = window
		mergedConfig.Vars = vars
		mergedConfig.Decay = decay
		return mergedConfig, nil
	}
}

// remove any duplicate Tuples, or Tuples that are less than or equal
// to other Tuples, since they'll get tracked automatically
func simplifySums(sums SumsConfig) SumsConfig {
	tupleMap := map[uint64]bool{}
	newSums := SumsConfig{}
	for i, m := range sums {
		// remove dupes
		if _, ok := tupleMap[m.hash()]; ok {
			continue
		}

		// if Tuple is less than or equal to an already existing Tuple,
		// skip this one
		leq := false
		for j, n := range sums {
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
		newSums = append(newSums, m)
	}

	return newSums
}

func validateConfig(config *CoreConfig) error {
	if config.Window == nil {
		return errors.New("config Window is not set")
	} else if *config.Window < 0 {
		return errors.Errorf("config has a negative window of %d", *config.Window)
	}

	if config.Vars != nil && *config.Vars < 2 {
		return errors.Errorf("config has less than 2 vars: %d < 2", *config.Vars)
	}

	if config.Decay != nil {
		if *config.Decay <= 0 || *config.Decay >= 1 {
			return errors.Errorf("config has a decay of %f, which is not in (0, 1)", *config.Decay)
		} else if *config.Window > 0 {
			return errors.New("config cannot have Decay set with a nonzero window")
		}
	}

	for _, tuple := range config.Sums {
		err := validateTuple(tuple, config)
		if err != nil {
			return err
		}
	}

	if config.Vars == nil && len(config.Sums) == 0 {
		return errors.New("config Vars is not set and cannot be inferred from empty Sums")
	}

	return nil
}

func validateTuple(tuple Tuple, config *CoreConfig) error {
	if len(tuple) != len(config.Sums[0]) {
		return errors.New("sums have differing length")
	} else if len(tuple) < 2 {
		return errors.Errorf("config has a Tuple (%v) with length %d < 2", tuple, len(tuple))
	}

	if config.Vars != nil && len(tuple) != *config.Vars {
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

	if config.Decay == nil {
		config.Decay = defaultConfig.Decay
	}

	return config
}
