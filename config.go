package stream

import "github.com/pkg/errors"

// CoreConfig is the struct containing configuration options for
// instantiating a Core object.
type CoreConfig struct {
	Sums        SumsConfig
	Window      *int
	PushMetrics []Metric
}

var defaultConfig = &CoreConfig{
	Sums:        map[int]bool{1: true},
	Window:      IntPtr(0),
	PushMetrics: nil,
}

// SumsConfig is an alias for a map of ints to bools; this configures
// the sums that a Core object will track.
type SumsConfig map[int]bool

func (s1 SumsConfig) add(s2 SumsConfig) {
	for k := range s2 {
		s1[k] = true
	}
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
			Sums:        SumsConfig{},
			PushMetrics: []Metric{},
		}

		for _, config := range configs {
			if config.Sums != nil {
				mergedConfig.Sums.add(config.Sums)
			}

			if config.Window != nil {
				if window == nil {
					window = IntPtr(*config.Window)
				} else if *window != *config.Window {
					return nil, errors.New("configs have differing windows")
				}
			}

			if config.PushMetrics != nil {
				mergedConfig.PushMetrics = append(mergedConfig.PushMetrics, config.PushMetrics...)
			}
		}

		mergedConfig.Window = window
		return mergedConfig, nil
	}
}

func validateConfig(config *CoreConfig) error {
	if config.Window != nil && *config.Window < 0 {
		return errors.New("config window is negative")
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

	// default Push is nil, no need to set

	return config
}
