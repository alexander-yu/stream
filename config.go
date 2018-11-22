package stream

import "github.com/pkg/errors"

// CoreConfig is the struct containing configuration options for
// instantiating a Core object.
type CoreConfig struct {
	Sums   SumsConfig
	Window *int
}

var defaultConfig = &CoreConfig{
	Sums:   map[int]bool{1: true},
	Window: IntPtr(1),
}

// SumsConfig is an alias for a map of ints to bools; this configures
// the sums that a Core object will track.
type SumsConfig map[int]bool

func (s1 SumsConfig) add(s2 SumsConfig) {
	for k := range s2 {
		s1[k] = true
	}
}

// MergeConfigs merges a slice of CoreConfig objects.
func MergeConfigs(configs []*CoreConfig) (*CoreConfig, error) {
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
					window = IntPtr(*config.Window)
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
	if config.Sums != nil && len(config.Sums) == 0 {
		errors.New("config sums map is not nil but empty")
	}

	if config.Window != nil && *config.Window <= 0 {
		errors.New("config window is nonpositive")
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
