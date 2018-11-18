package stream

import "github.com/pkg/errors"

// StatsConfig is the struct containing configuration options for
// instantiating a Stats object.
type StatsConfig struct {
	Sums   map[int]bool
	Window *int
	Median *bool
}

var defaultConfig = &StatsConfig{
	Sums:   map[int]bool{1: true},
	Window: IntPtr(1),
	Median: BoolPtr(true),
}

func validateConfig(config *StatsConfig) error {
	if config.Sums != nil && len(config.Sums) == 0 {
		errors.New("config sums map is empty")
	}

	if config.Window != nil && *config.Window <= 0 {
		errors.New("config window is nonpositive")
	}

	return nil
}

func setConfigDefaults(config *StatsConfig) *StatsConfig {
	if config.Sums == nil {
		config.Sums = defaultConfig.Sums
	}

	if config.Window == nil {
		config.Window = defaultConfig.Window
	}

	if config.Median == nil {
		config.Median = defaultConfig.Median
	}

	return config
}
