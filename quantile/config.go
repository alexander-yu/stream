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
	Impl          *Impl
}

var defaultConfig = &Config{
	Window:        stream.IntPtr(0),
	Interpolation: Linear.Ptr(),
	Impl:          AVL.Ptr(),
}

func validateConfig(config *Config) error {
	if config.Window == nil {
		return errors.New("config Window is not set")
	} else if config.Interpolation == nil {
		return errors.New("config Interpolation is not set")
	} else if config.Impl == nil {
		return errors.New("config Impl is not set")
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

	if !config.Impl.Valid() {
		return errors.Errorf(
			"config has an invalid Impl of %d",
			*config.Impl,
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

	if config.Impl == nil {
		config.Impl = defaultConfig.Impl
	}

	return config
}
