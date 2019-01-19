package quantile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexander-yu/stream"
)

func TestValidateConfig(t *testing.T) {
	t.Run("fail: config without Window is invalid", func(t *testing.T) {
		config := &Config{
			Interpolation: Linear.Ptr(),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Window is not set")
	})

	t.Run("fail: config without Interpolation is invalid", func(t *testing.T) {
		config := &Config{
			Window: stream.IntPtr(0),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Interpolation is not set")
	})

	t.Run("fail: config without Impl is invalid", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(0),
			Interpolation: Linear.Ptr(),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, "config Impl is not set")
	})

	t.Run("fail: config with a negative window is invalid", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(-1),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf("config has a negative window of %d", -1))
	})

	t.Run("fail: config with unsupported Interpolation is invalid", func(t *testing.T) {
		interpolation := Interpolation(-1)
		config := &Config{
			Window:        stream.IntPtr(0),
			Interpolation: &interpolation,
			Impl:          AVL.Ptr(),
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf(
			"config has an invalid Interpolation of %d",
			*config.Interpolation,
		))
	})

	t.Run("fail: config with unsupported Impl is invalid", func(t *testing.T) {
		impl := Impl(-1)
		config := &Config{
			Window:        stream.IntPtr(0),
			Interpolation: Linear.Ptr(),
			Impl:          &impl,
		}
		err := validateConfig(config)
		assert.EqualError(t, err, fmt.Sprintf(
			"config has an invalid Impl of %d",
			*config.Impl,
		))
	})

	t.Run("pass: valid config is valid", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(0),
			Interpolation: Linear.Ptr(),
			Impl:          AVL.Ptr(),
		}
		err := validateConfig(config)
		assert.NoError(t, err)
	})
}

func TestSetConfigDefaults(t *testing.T) {
	t.Run("pass: empty config is overridden", func(t *testing.T) {
		config := &Config{}
		config = setConfigDefaults(config)

		expectedConfig := defaultConfig

		assert.Equal(t, expectedConfig, config)
	})

	t.Run("pass: provided fields are kept", func(t *testing.T) {
		config := &Config{
			Window:        stream.IntPtr(2),
			Interpolation: Midpoint.Ptr(),
			Impl:          AVL.Ptr(),
		}
		config = setConfigDefaults(config)

		expectedConfig := &Config{
			Window:        stream.IntPtr(2),
			Interpolation: Midpoint.Ptr(),
			Impl:          AVL.Ptr(),
		}

		assert.Equal(t, expectedConfig, config)
	})
}
